
# GO Library to use JUJU commands to bring up a Kubernetes cluster

### Pre-reqs

`go get gopkg.in/yaml.v2`

## Usage

With a valid `manifest.yaml` to pass along creds and cloud info, this can be used for :

- `Spinup()` - will spinup cluster from specified creds/cloud/bundle
- `ClusterReady()` - will return boolean corresponding to readiness of cluster
- `DisplayStatus()` - will display results of running `juju status`
- `GetKubeConfig()` - will print out kubeconfig to stdout
- `DestroyCluster()` - will tear down juju controller and associated cluster

## Notes:

`manifest.yaml` will for include credentials and cloud information to be consumed by gogo.go


Should be of the format:
```
credentials:
  aws:
    <name>:
      auth-type: access-key
      access-key: <aws-access-key>
      secret-key: <aws-secret-key>
  lab:
    <name>:
      auth-type: oauth1
      maas-oauth: <maas-api-key>

clouds:
   lab:
      type: maas
      auth-types: [oauth1]
      endpoint: <your-maas-url>
```

Sample file that would use this library:

```package main

import (
	"github.com/dstorck/gogo"
)

var testRun = gogo.Juju{
	Name:     "test-cluster",
	Bundle:   "cs:bundle/kubernetes-core-306",
	Manifest: "manifest.yaml",
}

// current available commands, not meant to be run all at once
func main() {
	 testRun.Spinup()
   testRun.DisplayStatus()
	 testRun.ClusterReady()
   testRun.GetKubeConfig()
	 testRun.DestroyCluster()
}
```
