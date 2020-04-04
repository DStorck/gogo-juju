# GO Library to use JUJU commands to bring up a Kubernetes cluster
[![Go Report Card](https://goreportcard.com/badge/github.com/DStorck/gogo)](https://goreportcard.com/report/github.com/DStorck/gogo)
### Pre-reqs
`go get gopkg.in/yaml.v2`

## Usage
Currently can bring up kubernetes cluster in CNCT maas lab or on aws.

You will need to set the instance of `Juju` with the following fields:

### JujuDataPrefix variable
In order to run the juju cli in discrete environments, a different JUJU_DATA path is set per
Juju.Name. This path prefix defaults to `/tmp` and will result in the loss of the juju cli
state as /tmp is ephemeral.

Set this variable to a path with persistant storage:
```go
import (
    "github.com/dstorck/gogo"
)

gogo.JujuDataPrefix = "/data"
```

### Juju struct options
* Note - While several fields are optional, you must include either `MaasCl` __and__ `MaasCr` __or__ `AwsCl` __and__ `AwsCr`

| Field Name     | Required    | Type            | Description                                 |
| -------------- | ----------- | --------------- | ------------------------------------------- |
| Kind           | __Required__| CloudKind       | must use one of the const: `Maas` or `Aws`  |
| Name           | __Required__| String          | JUJU_DATA path and juju clustername for MAAS|
| Bundle         | __Required__| String          | ex "cs:bundle/canonical-kubernetes-193"     |
| p              |  Optional   | Parallel        | used for multiple cluster creation          |
| MaasCl         | Optional    | MaasCloud       | maas cloud details                          |
| MaasCr         | Optional    | MaasCredentials | maas credential details                     |
| AwsCl          | Optional    | AwsCloud        | aws cloud details                           |
| AwsCr          | Optional    | AwsCredentials  | aws credential details                      |

## MaasCl Options
| Field Name     | Required    | Type            | Description                             |
| -------------- | ----------- | --------------- | --------------------------------------- |
| Endpoint       |__Required__ | String          | maas url ex-"http://your-ip/MAAS/api/2.0"  |

## MaasCr Options
| Field Name     | Required    | Type            | Description                             |
| -------------- | ----------- | --------------- | --------------------------------------- |
| Username       | __Required__| String          | maas username 													 |
| MaasOauth      | __Required__| String          | maas api key                            |

## AwsCl Options
| Field Name     | Required    | Type            | Description                             |
| -------------- | ----------- | --------------- | --------------------------------------- |
| Region         | Optional    | String          | ex "aws/us-west-2"                      |

## AwsCr Options
| Field Name     | Required    | Type            | Description                             |
| -------------- | ----------- | --------------- | --------------------------------------- |
| Username       | __Required__| String          | your aws username                       |
| AccessKey      | __Required__| String          | aws access key                          |
| SecretKey      | __Required__| String          | aws secret key                          |

### Working commands:

- `SetAWSCreds()` - sets aws credentials for use with juju
- `SetMAASCreds()` sets maas credentials for use with juju
- `SetMAASCloud()` - sets maas cloud information for use with juju
- `Spinup()` - spins up cluster from specified creds/cloud/bundle (includes setting cloud and credentials)
- `ControllerReady()` - returns boolean with status of controller availability
- `GetStatus()` - returns result of running juju status
- `ClusterReady()` - returns boolean corresponding to readiness of cluster
- `GetKubeConfig()` - returns the contents of the kubeconfig file
- `DestroyCluster()` - tears down juju controller and associated cluster
- `DestroyComplete()` returns boolean corresponding to successful destruction of cluster


### Sample file to bring up maas cluster on Samsung's CNCT nuc lab:

```go
package main

import (
	"fmt"

	"github.com/dstorck/gogo"
)

var testRun = gogo.Juju{
	Kind:   gogo.Maas,
	Name:   "test-cluster",
	Bundle: "cs:bundle/kubernetes-core-306",
	MaasCl: myMaasCloud,
	MaasCr: myMaasCreds,
}

var myMaasCloud = gogo.MaasCloud{
	Endpoint: "<your-maas-url>",
}
var myMaasCreds = gogo.MaasCredentials{
	Username:  "<username>",
	MaasOauth: "<maas-api-key>",
}

// current available commands, not meant to be run all at once
func main() {
  testRun.SetMAASCreds()
  testRun.SetMAASCloud()
  testRun.Spinup()
  status, _ := testRun.GetStatus()
	fmt.Println(status)
  testRun.ClusterReady()
  config,_ := testRun.GetKubeConfig()
	fmt.Println(string(config))
  testRun.DestroyCluster()
  testRun.DestroyComplete()
}
```

### Sample file to bring up aws cluster
```go
package main

import (
	"github.com/dstorck/gogo"
)

var testRun = gogo.Juju{
	Kind:   gogo.Aws,
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
  // testRun.GetStatus()
  // testRun.Spinup()
  // testRun.ControllerReady()
  // testRun.ClusterReady()
  // testRun.GetKubeConfig()
  // testRun.DestroyCluster()
  // testRun.DestroyComplete()
}
```
