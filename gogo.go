package gogo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var jStats jujuStatus

// JujuDataPrefix is the path prefix used for the JUJU_DATA environment variable
// this path will store required juju state and should be persistent
var JujuDataPrefix = "/tmp/"

// Spinup will create one cluster
func (j *Juju) Spinup() {
	controller := ""
	user := ""
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	if j.Kind == Aws {
		j.SetAWSCreds()
		controller = j.AwsCl.Region
		user = j.AwsCr.Username
	} else if j.Kind == Maas {
		j.SetMAASCloud()
		j.SetMAASCreds()
		controller = j.MaasCl.Type
		user = j.MaasCr.Username
	}

	credscommand := "--credential=" + user

	cmd := exec.Command("juju", "bootstrap", controller, credscommand)

	// cmd := exec.Command("juju", "bootstrap", controller) // with aws this is is expecting region ex - juju bootstrap aws/us-west-2
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "bootstrap")

	cmd = exec.Command("juju", "add-model", j.Name, credscommand)
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.CombinedOutput()
	commandResult(out, err, "add-model")

	cmd = exec.Command("juju", "deploy", j.Bundle)
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.CombinedOutput()
	commandResult(out, err, "deploy")
}

// DisplayStatus will ask juju for status
func (j *Juju) DisplayStatus() {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "status")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "display status")
}

// ClusterReady will check status and return true if cluster is running
func (j *Juju) ClusterReady() bool {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "status", "--format=json")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("ClusterReady() failed with %s\n", err)
	}

	json.Unmarshal([]byte(out), &jStats)

	for k := range jStats.Machines {
		machineStatus := jStats.Machines[k].MachStatus["current"]
		if machineStatus != "started" {
			fmt.Println("Cluster Not Ready")
			return false
		}
	}

	for k := range jStats.ApplicationResults {
		appStatus := jStats.ApplicationResults[k].AppStatus["current"]
		if appStatus != "active" {
			fmt.Println("Cluster Not Ready")
			return false
		}
	}

	fmt.Println("Cluster Ready")
	return true
}

// GetKubeConfig returns the kubeconfig file contents
func (j *Juju) GetKubeConfig() ([]byte, error) {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "ssh", "kubernetes-master/0", "cat", "config")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("GetKubeConfig failed: %s", err)
	}
	return out, nil
}

// DestroyCluster will kill off one cluster
func (j *Juju) DestroyCluster() {
	controller := ""
	if j.Kind == Aws {
		controller = j.AwsCl.Region
	} else if j.Kind == Maas {
		controller = j.MaasCl.Type
	}
	controller = strings.Replace(controller, "/", "-", -1)

	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name
	cmd := exec.Command("juju", "destroy-controller", "--destroy-all-models", controller, "-y")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "destroy-controller")
}

func commandResult(out []byte, err error, command string) {
	fmt.Printf("\n%s\n", string(out))
	if err != nil {
		log.Fatalf("%s failed with %s\n", command, err)
	}
}

// Create is an example of spinning up multiple clusters
func (j *Juju) Create(clusters []string) {
	// clusters := []string{"d8048274-2bc6-49bf-81fd-846aeaddf2fe", "97c19eda-7aeb-4eee-a35c-57dc3755d98f"}

	// for _, cluster := range clusters {
	// 	j.p.wg.Add(1)
	// 	go j.Spinup()
	// }
	// j.p.wg.Wait()
}
