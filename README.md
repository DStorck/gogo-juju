# GO Library to use JUJU commands to bring up a Kubernetes cluster

### Pre-reqs
`go get gopkg.in/yaml.v2`

## Usage
Currently can bring up kubernetes cluster in CNCT maas lab or on aws.

### Working commands:

- `SetAWSCreds()` - sets aws credentials for use with juju
- `SetMAASCreds()` sets maas credentials for use with juju
- `SetMAASCloud()` - sets maas cloud information for use with juju
- `Spinup()` - will spinup cluster from specified creds/cloud/bundle
- `DisplayStatus()` - will display results of running juju status
- `ClusterReady()` - will return boolean corresponding to readiness of cluster
- `GetKubeConfig()` - will print out kubeconfig to stdout
- `DestroyCluster()` - will tear down juju controller and associated cluster


### Sample file to bring up maas cluster on Samsung's CNCT nuc lab:

```
package main

import (
	"github.com/dstorck/gogo"
)

var testRun = gogo.Juju{
	Kind:   "maas",
	Name:   "test-cluster",
	Bundle: "cs:bundle/kubernetes-core-306",
	MaasCl: myMaasCloud,
	MaasCr: myMaasCreds,
}

var myMaasCloud = gogo.Cloud{
	Type:     "lab",
	Endpoint: "<your-maas-url>",
}
var myMaasCreds = gogo.Credentials{
	CloudName: "lab",
	Username:  "deirdre",
	MaasOauth: "<maas-password>",
}

// current available commands, not meant to be run all at once
func main() {
  testRun.SetMAASCreds()
  testRun.SetMAASCloud()
  testRun.Spinup()
  testRun.DisplayStatus()
  testRun.ClusterReady()
  testRun.GetKubeConfig()
  testRun.DestroyCluster()
}
```

### Sample file to bring up aws cluster
```
package main

import (
	"github.com/dstorck/gogo"
)

var testRun = gogo.Juju{
	Kind:   "aws",
	Name:   "test-cluster",
	Bundle: "cs:bundle/kubernetes-core-306",
	AwsCr:  myAWScreds,
	AwsCl:  myAWScloud,
}

var myAWScreds = gogo.AWSCredentials{
	Username:  "<your-aws-username>",
	AccessKey: "<your-access-key>",
	SecretKey: "<your-secret-key",
}

var myAWScloud = gogo.AWSCloud{
	Region: "aws/us-west-2",
}

// currently available commands, not meant to be run all at once
func main() {
  // testRun.SetAWSCreds()
  // testRun.DisplayStatus()
  // testRun.Spinup()
  // testRun.ClusterReady()
  // testRun.GetKubeConfig()
  // testRun.DestroyCluster()
}
```
