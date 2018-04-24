package gogo

import (
	"errors"

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
func CreateMAASCredsYaml(name string, username string, maasOauth string) (string, error) {
	if name == "" {
		return "", errors.New("Name must not be empty")
	}
	if username == "" {
		return "", errors.New("User must not be empty")
	}
	if maasOauth == "" {
		return "", errors.New("Maas-Oauth must not be empty")
	}
	lab := newSuperCred(creds{
		name: user{
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

// credentials:
//   lab:
//     <username>:
//       auth-type: oauth1
//       maas-oauth: <your-maas-secret>
