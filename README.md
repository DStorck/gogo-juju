
# GO Library to use JUJU commands to bring up a Kubernetes cluster

### Pre-reqs

`go get gopkg.in/yaml.v2`

## Usage

With a valid `manifest.yaml` to pass along creds and cloud info, this can be used for :

- `SetCloudAndCreds()` - sets cloud and credentials information for use with `juju`
- `Spinup()` - will spinup cluster from specified creds/cloud/bundle
- `DisplayStatus()` - will display results of running `juju status`
- `ClusterReady()` - will return boolean corresponding to readiness of cluster
- `GetKubeConfig()` - will print out kubeconfig to stdout
- `DestroyCluster()` - will tear down juju controller and associated cluster

## Notes:

We will pass in cloud and credential information from the Custom Resource.

Sample file that would use this library:

```package main

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
	Endpoint: "http://192.168.2.24/MAAS/api/2.0",
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
