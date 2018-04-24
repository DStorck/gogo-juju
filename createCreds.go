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

type creds map[string]User

type User map[string]Auth

type Auth struct {
	AuthType  string `yaml:"auth-type,flow"`
	MaasOauth string
}

// CreateMAASCloudYaml is used to create the yaml string to pass to "juju add-cloud"
func CreateMAASCredsYaml(name string, user, string, maasOauth string) (string, error) {
	if name == "" {
		return "", errors.New("Name must not be empty")
	}
	if user == "" {
		return "", errors.New("User must not be empty")
	}
	if maasOauth == "" {
		return "", errors.New("Maas-Oauth must not be empty")
	}
	lab := newSuperCred(creds{
		name: User{
			user: Auth{
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

// func main() {
// 	out, err := CreateMAASCloudYaml("lab", "http://192.168.2.24/MAAS/api/2.0")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(out)
// }

// credentials:
//   lab:
//     deirdre:
//       auth-type: oauth1
//       maas-oauth: 85uf5sEpqyNHy6ALVy:fpEbqGPz9tS9qfJNxc:8AfrxLgLKyTPUC4679jkZtMq7GhG4UwJ

// clouds:
//   lab:
//     type: maas
//     auth-types: [oauth1]
//     endpoint: http://192.168.2.24/MAAS/api/2.0/
