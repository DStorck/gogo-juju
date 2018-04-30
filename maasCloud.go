package gogo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type superCloud map[string]clouds

func newSuperCloud(c clouds) *superCloud {
	sc := superCloud{"clouds": c}
	return &sc
}

type clouds map[string]cloud

type cloud struct {
	Type      string
	AuthTypes []string `yaml:"auth-types,flow"`
	Endpoint  string
}

// CreateMAASCloudYaml is used to create the yaml string to pass to "juju add-cloud"
func CreateMAASCloudYaml(name string, endpoint string) (string, error) {
	if name == "" {
		return "", errors.New("Name must not be empty")
	}
	if endpoint == "" {
		return "", errors.New("Endpoint must not be empty")
	}
	lab := newSuperCloud(clouds{
		name: {
			Type:      "maas",
			AuthTypes: []string{"oauth1"},
			Endpoint:  endpoint,
		},
	})

	output, err := yaml.Marshal(lab)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// SetMAASCloud will run juju add-cloud with maasCloud yaml created above
func (j *Juju) SetMAASCloud() {
	tmp := "JUJU_DATA=/tmp/" + j.Name

	cloudInfo, err := CreateMAASCloudYaml(j.MaasCl.Type, j.MaasCl.Endpoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cloudInfo)
	cmd := exec.Command("juju", "add-cloud", j.MaasCl.Type, "-f", "/dev/stdin", "--replace")
	cmd.Stdin = strings.NewReader(cloudInfo)
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "add-cloud")
}
