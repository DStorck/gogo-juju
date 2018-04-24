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

// SetCloudAndCreds will grab cloud and credential information and set it
func (j *Juju) SetCloudAndCreds() {
	tmp := "JUJU_DATA=/tmp/" + j.Name

	manifest, err := CreateMAASCloudYaml(j.Cl.Type, j.Cl.Endpoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(manifest)
	cmd := exec.Command("juju", "add-cloud", "lab", "-f", "/dev/stdin", "--replace")
	cmd.Stdin = strings.NewReader(manifest)
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "add-cloud")

	creds, err := CreateMAASCredsYaml(j.Cr.CloudName, j.Cr.Username, j.Cr.MaasOauth)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(creds)

	cmd = exec.Command("juju", "add-credential", "lab", "-f", "/dev/stdin", "--replace")
	cmd.Stdin = strings.NewReader(creds)
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.CombinedOutput()
	commandResult(out, err, "add-credential")
}

// Spinup will create one cluster
func (j *Juju) Spinup() {
	tmp := "JUJU_DATA=/tmp/" + j.Name

	j.SetCloudAndCreds()

	cmd := exec.Command("juju", "bootstrap", j.Cl.Type)
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "bootstrap")

	cmd = exec.Command("juju", "deploy", j.Bundle)
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.CombinedOutput()
	commandResult(out, err, "deploy")
}

// DisplayStatus will ask juju for status
func (j *Juju) DisplayStatus() {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "status")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "display status")
}

// ClusterReady will check status and return true if cluster is running
func (j *Juju) ClusterReady() bool {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "status", "--format=json")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s failed with %s\n", "get cluster deets", err)
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

// GetKubeConfig will cat out kubernetes config to stdout
func (j *Juju) GetKubeConfig() {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "ssh", "kubernetes-master/0", "cat", "config")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "get kube config")
}

// DestroyCluster will kill off one cluster
func (j *Juju) DestroyCluster() {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "destroy-controller", "--destroy-all-models", "lab", "-y")
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
