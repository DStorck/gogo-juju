package gogo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type superCred map[string]creds

func newSuperCred(c creds) *superCred {
	sc := superCred{"credentials": c}
	return &sc
}

type creds map[string]user

type user map[string]auth

type auth struct {
	AuthType  string `yaml:"auth-type"`
	MaasOauth string `yaml:"maas-oauth"`
}

// CreateMAASCredsYaml is used to create the yaml string to pass to "juju add-credential"
func CreateMAASCredsYaml(cloudName string, username string, maasOauth string) (string, error) {
	if cloudName == "" {
		return "", errors.New("cloudName must not be empty")
	}
	if username == "" {
		return "", errors.New("User must not be empty")
	}
	if maasOauth == "" {
		return "", errors.New("Maas-Oauth must not be empty")
	}
	lab := newSuperCred(creds{
		cloudName: user{
			username: auth{
				AuthType:  "oauth1",
				MaasOauth: maasOauth,
			},
		},
	})

	output, err := yaml.Marshal(lab)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// SetMAASCreds will pass in maas credentials to juju add-credential
func (j *Juju) SetMAASCreds() {
	tmp := "JUJU_DATA=/tmp/" + j.Name

	creds, err := CreateMAASCredsYaml(j.MaasCr.CloudName, j.MaasCr.Username, j.MaasCr.MaasOauth)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(creds)

	cmd := exec.Command("juju", "add-credential", j.MaasCl.Type, "-f", "/dev/stdin", "--replace")
	cmd.Stdin = strings.NewReader(creds)
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "add-credential")
}
