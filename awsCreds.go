package gogo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type superAWSCred map[string]awsCreds

func newSuperAWSCred(c awsCreds) *superAWSCred {
	sc := superAWSCred{"credentials": c}
	return &sc
}

type awsCreds map[string]awsUser

type awsUser map[string]awsAuth

type awsAuth struct {
	AuthType  string `yaml:"auth-type"`
	AccessKey string `yaml:"access-key"`
	SecretKey string `yaml:"secret-key"`
}

// CreateAWSCredsYaml is used to create the yaml string to pass to "juju add-credential"
func CreateAWSCredsYaml(username string, accessKey string, secretKey string) (string, error) {
	if username == "" {
		return "", errors.New("User must not be empty")
	}
	if accessKey == "" {
		return "", errors.New("Access Key must not be empty")
	}
	if secretKey == "" {
		return "", errors.New("Secret Key must not be empty")
	}
	aws := newSuperAWSCred(awsCreds{
		"aws": awsUser{
			username: awsAuth{
				AuthType:  "access-key",
				AccessKey: accessKey,
				SecretKey: secretKey,
			},
		},
	})

	output, err := yaml.Marshal(aws)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// SetAWSCreds will grab and credential information and set it
func (j *Juju) SetAWSCreds() error {
	tmp := "JUJU_DATA=" + JujuDataPrefix + j.Name

	creds, err := CreateAWSCredsYaml(j.AwsCr.Username, j.AwsCr.AccessKey, j.AwsCr.SecretKey)
	if err != nil {
		return fmt.Errorf("SetAWSCreds error: %s", err)
	}
	fmt.Println(creds)

	cmd := exec.Command("juju", "add-credential", "aws", "-f", "/dev/stdin", "--replace")
	cmd.Stdin = strings.NewReader(creds)
	cmd.Env = append(os.Environ(), tmp)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("SetAWSCreds error: %s", err)
	}
	return nil
}
