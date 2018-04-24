package gogo

import (
	"errors"

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

// func main() {
// 	out, err := CreateMAASCloudYaml("lab", "http://192.168.2.24/MAAS/api/2.0")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(out)
// }

// clouds:
//   lab:
//     type: maas
//     auth-types: [oauth1]
//     endpoint: http://192.168.2.24/MAAS/api/2.0/
