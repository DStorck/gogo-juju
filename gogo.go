package gogo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

// Juju struct - name , bundle, manifest, and connection to type Parallel
type Juju struct {
	name     string
	bundle   string
	manifest string
	p        Parallel
}

// Parallel sets the waitgroup
type Parallel struct {
	wg sync.WaitGroup
}

// DisplayStatus will ask juju for status
func (j *Juju) DisplayStatus() {
	cmd := exec.Command("juju", "status")
	cmd.Env = append(os.Environ(), "JUJU_DATA=/tmp/"+j.name)
	out, err := cmd.CombinedOutput()
	fmt.Printf("\n%s\n", string(out))
	if err != nil {
		log.Fatalf("add-cloud failed with %s\n", err)
	}
}

// Spinup will create one cluster
func (j *Juju) Spinup() {
	cmd := exec.Command("juju", "add-cloud", "lab", "-f", j.manifest, "--replace")
	cmd.Env = append(os.Environ(), "JUJU_DATA=/tmp/"+j.name)
	out, err := cmd.CombinedOutput()
	fmt.Printf("\n%s\n", string(out))
	if err != nil {
		log.Fatalf("add-cloud failed with %s\n", err)
	}
	cmd = exec.Command("juju", "add-credential", "lab", "-f", j.manifest, "--replace")
	cmd.Env = append(os.Environ(), "JUJU_DATA=/tmp/"+j.name)
	out, err = cmd.CombinedOutput()
	fmt.Printf("\n%s\n", string(out))
	if err != nil {
		log.Fatalf("add-credential failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	cmd = exec.Command("juju", "bootstrap", "lab")
	cmd.Env = append(os.Environ(), "JUJU_DATA=/tmp/"+j.name)
	out, err = cmd.CombinedOutput()
	fmt.Printf("\n%s\n", string(out))
	if err != nil {
		log.Fatalf("Boostrap controller failed with %s\n", err)
	}
	cmd = exec.Command("juju", "deploy", j.bundle)
	cmd.Env = append(os.Environ(), "JUJU_DATA=/tmp/"+j.name)
	out, err = cmd.CombinedOutput()
	fmt.Printf("\n%s\n", string(out))
	if err != nil {
		log.Fatalf("Deploy failed with %s\n", err)
	}
	j.p.wg.Done()
}

// DestroyCluster will kill off one cluster
func (j *Juju) DestroyCluster(cluster string) {
	cmd := exec.Command("juju", "destroy-controller", "-y", "lab")
	// cmd.Stdin = strings.NewReader("y")
	cmd.Env = append(os.Environ(), "JUJU_DATA=/tmp/"+j.name)
	out, err := cmd.CombinedOutput()
	fmt.Printf("\n%s\n", string(out))
	if err != nil {
		log.Fatalf("Destroy failed with %s\n", err)
	}
}

var clusters = []string{"deirdre-test"}

// Create will create all clusters in an array given their names
func (j *Juju) Create(clusters []string) {
	// clusters := []string{"d8048274-2bc6-49bf-81fd-846aeaddf2fe", "97c19eda-7aeb-4eee-a35c-57dc3755d98f"}

	// for _, cluster := range clusters {
	// 	j.p.wg.Add(1)
	// 	go j.Spinup()
	// }
	// j.p.wg.Wait()
}
