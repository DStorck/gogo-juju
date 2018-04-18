package gogo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

// Juju struct - name , bundle, manifest, and connection to type Parallel
// manifest is name of file to contain credential and cloud details
type Juju struct {
	Name     string
	Bundle   string
	Manifest string
	p        Parallel
}

// Parallel sets the waitgroup
type Parallel struct {
	wg sync.WaitGroup
}

// unit status is coming from applications.<app-name>.units.<unit-name>.workload-status.current - should be active
// app status applications.<app-name>.application-status.current - should be active
// machine status - machines.<machine-id>.juju-status.current - should be started

var jsonResults map[string]interface{}

// GetClusterDeets blah
func (j *Juju) GetClusterDeets() {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "status", "--format=json")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s failed with %s\n", "get deets", err)
	}
	json.Unmarshal([]byte(out), &jsonResults)
	appNames := jsonResults["applications"].(map[string]interface{})
	for key, _ := range appNames {
		// Each value is an interface{} type, that is type asserted as a string
		fmt.Println(key)
	}

}

// Spinup will create one cluster
func (j *Juju) Spinup() {
	tmp := "JUJU_DATA=/tmp/" + j.Name

	cmd := exec.Command("juju", "add-cloud", "lab", "-f", j.Manifest, "--replace")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "add-cloud")

	cmd = exec.Command("juju", "add-credential", "lab", "-f", j.Manifest, "--replace")
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.CombinedOutput()
	commandResult(out, err, "add-credential")

	cmd = exec.Command("juju", "bootstrap", "lab")
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.CombinedOutput()
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

// ClusterReady will block until all units, machines, and apps are displaying ready
func (j *Juju) ClusterReady() {
	fmt.Println("is your cluster ready? we shall see . ")
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

// Create will create all clusters in an array given their names
func (j *Juju) Create(clusters []string) {
	// clusters := []string{"d8048274-2bc6-49bf-81fd-846aeaddf2fe", "97c19eda-7aeb-4eee-a35c-57dc3755d98f"}

	// for _, cluster := range clusters {
	// 	j.p.wg.Add(1)
	// 	go j.Spinup()
	// }
	// j.p.wg.Wait()
}
