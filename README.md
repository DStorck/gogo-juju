# GO Library to use JUJU commands to bring up a Kubernetes cluster
### Pre-reqs
`go get gopkg.in/yaml.v2`

## Usage
Currently can bring up kubernetes cluster in CNCT maas lab. Check branch aws-cluster for use with aws.

## Working commands:

- `SetCloudAndCreds()` - sets cloud and credentials information for use with juju
- `Spinup()` - will spinup cluster from specified creds/cloud/bundle
- `DisplayStatus()` - will display results of running juju status
- `ClusterReady()` - will return boolean corresponding to readiness of cluster
- `GetKubeConfig()` - will print out kubeconfig to stdout
- `DestroyCluster()` - will tear down juju controller and associated cluster

Sample file that would use this library:
```
import (
	"github.com/dstorck/gogo"
)

var testRun = gogo.Juju{
	Name:     "test-cluster",
	Bundle:   "cs:bundle/kubernetes-core-306",
	Cl:       myCloud,
	Cr:       myCreds,
}

var myCloud = gogo.Cloud{
	Type:     "lab",
	Endpoint: "<maas-url>",
}

var myCreds = gogo.Credentials{
	CloudName: "lab",
	Username:  "<your-maas-username>",
	MaasOauth: "<your-maas-secret>",
}


// current available commands, not meant to be run all at once
func main() {
  testRun.SetCloudAndCreds()
  testRun.Spinup()
  testRun.DisplayStatus()
  testRun.ClusterReady()
  testRun.GetKubeConfig()
  testRun.DestroyCluster()
}
```
