
# GO Library to use JUJU commands to bring up a Kubernetes cluster

With a valid `manifest.yaml` to pass along creds and cloud info, this can be used for :

- `Spinup()`
- `DisplayStatus()`
- `ClusterReady()`
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
