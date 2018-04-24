
# GO Library to use JUJU commands to bring up a Kubernetes cluster

With a valid `manifest.yaml` to pass along creds and cloud info, this can be used for :

- `Spinup()`
- `GetJujuStatus()`
- `ClusterReady()` - WIP
- `GetKubeConfig()`
- `DestroyCluster()`

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
	 testRun.GetJujuStatus()
   testRun.GetKubeConfig()
	 testRun.DestroyCluster()
}
```
