package gogo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

var jStats jujuStatus
var jControllers jujuControllers
var jModels jujuModels

// JujuDataPrefix is the path prefix used for the JUJU_DATA environment variable
// this path will store required juju state and should be persistent
var JujuDataPrefix = "/tmp/"

var log = logrus.New()

// Spinup will create one cluster
func (j *Juju) Spinup() error {
	var controller string
	var user string
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	if j.Kind == Aws {
		err := j.SetAWSCreds()
		if err != nil {
			return fmt.Errorf("Spinup error with setAwsCreds: %s", err)
		}
		controller = j.AwsCl.Region
		user = j.AwsCr.Username
	} else if j.Kind == Maas {
		err := j.SetMAASCloud()
		if err != nil {
			return fmt.Errorf("Spinup error with setMaasCloud: %s", err)
		}
		err = j.SetMAASCreds()
		if err != nil {
			return fmt.Errorf("Spinup error with setMaasCreds: %s", err)
		}
		controller = j.Name
		user = j.MaasCr.Username
	} else {
		return errors.New("Spinup error: Juju.Kind must be a supported cloud")
	}

	credscommand := "--credential=" + user

	cmd := exec.Command("juju", "bootstrap", controller, credscommand)
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.Output()
	if err != nil {
		// return fmt.Errorf("Spinup error with bootstrap: %s", err)
		return fmt.Errorf("Spinup error with bootstrap: %v: %s", err, err.(*exec.ExitError).Stderr)
	}
	log.Debug(string(out))

	cmd = exec.Command("juju", "add-model", j.Name, credscommand)
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("Spinup error with add-model:%v: %s", err, err.(*exec.ExitError).Stderr)
	}
	log.Debug(string(out))

	cmd = exec.Command("juju", "deploy", j.Bundle)
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("Spinup error with deploy: %v: %s", err, err.(*exec.ExitError).Stderr)
	}
	log.Debug(string(out))

	return nil
}

// ControllerReady checks model status and returns bool
func (j *Juju) ControllerReady() (bool, error) {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "models", "--format=json")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("ControllerReady error: %v: %s", err, err.(*exec.ExitError).Stderr)
	}

	err = json.Unmarshal([]byte(out), &jModels)
	if err != nil {
		return false, fmt.Errorf("ControllerReady unmarshal error: %v: %s", err, err.(*exec.ExitError).Stderr)
	}

	log.Debugf("ControllerReady: %+v", jModels)
	for k := range jModels.Models {
		if jModels.Models[k].ShortName == j.Name {
			status := jModels.Models[k].Status["current"]
			if status == "available" {
				log.WithFields(logrus.Fields{"name": j.Name}).Info("Controller Ready")
				return true, nil
			}
		}
	}
	log.WithFields(logrus.Fields{"name": j.Name}).Info("Controller Not Ready")
	return false, nil
}

// GetStatus return juju status
func (j *Juju) GetStatus() (string, error) {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "status")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("GetStatus error: %v: %s", err, err.(*exec.ExitError).Stderr)
	}
	log.Debug(string(out))
	return string(out), nil
}

// ClusterReady will check status and return true if cluster is running
func (j *Juju) ClusterReady() (bool, error) {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "status", "--format=json")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("ClusterReady error: %v: %s", err, err.(*exec.ExitError).Stderr)
	}

	err = json.Unmarshal([]byte(out), &jStats)
	if err != nil {
		return false, fmt.Errorf("ClusterReady error: %v: %s", err, err.(*exec.ExitError).Stderr)
	}

	for k := range jStats.Machines {
		machineStatus := jStats.Machines[k].MachStatus["current"]
		if machineStatus != "started" {
			log.WithFields(logrus.Fields{"name": j.Name}).Info("Cluster Not Ready")
			return false, nil
		}
	}

	for k := range jStats.ApplicationResults {
		appStatus := jStats.ApplicationResults[k].AppStatus["current"]
		if appStatus != "active" {
			log.WithFields(logrus.Fields{"name": j.Name}).Info("Cluster Not Ready")
			return false, nil
		}
	}

	log.WithFields(logrus.Fields{"name": j.Name}).Info("Cluster Ready")
	return true, nil
}

// GetKubeConfig returns the kubeconfig file contents
func (j *Juju) GetKubeConfig() ([]byte, error) {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "ssh", "kubernetes-master/0", "cat", "config")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.Output()
	if err != nil {
		return []byte{}, fmt.Errorf("GetKubeConfig failed: %v: %s", err, err.(*exec.ExitError).Stderr)
	}
	return out, nil
}

// DestroyCluster will kill off one cluster
func (j *Juju) DestroyCluster() error {
	controller := ""
	if j.Kind == Aws {
		controller = j.AwsCl.Region
	} else if j.Kind == Maas {
		controller = j.Name
	} else {
		return errors.New("DestroyCluster: Juju.Kind must be a supported cloud")
	}
	controller = strings.Replace(controller, "/", "-", -1)

	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "destroy-controller", "--destroy-all-models", controller, "-y")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("DestroyCluster error: %v: %s", err, err.(*exec.ExitError).Stderr)
	}
	log.Debug(string(out))
	return nil
}

// DestroyComplete checks juju for controllers to make sure none are left
func (j *Juju) DestroyComplete() (bool, error) {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "controllers", "--format=json")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("DestroyComplete error: %v: %s", err, err.(*exec.ExitError).Stderr)
	}

	err = json.Unmarshal([]byte(out), &jControllers)
	if err != nil {
		return false, fmt.Errorf("DestroyComplete error: %s", err)
	}

	log.Debugf("DestroyComplete: %+v", jControllers)
	if len(jControllers.Controllers) < 1 {
		return true, nil
	}
	return false, nil
}
