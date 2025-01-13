# ci-secret-sync

There are secrets that are maintained outside the cluster that when changed, must be shadowed and dependent pods must be restarted. To faciliate this, `ci-secret-sync` checks for updates to secrets
and updates target secrets along with restarting pods. 

sample: 
```
sync:
  actions:
  - source:
      name: "pull-secret"
      namespace: "openshift-config"
      key: ".dockerconfigjson"
      type: "secret"
      schema: "REGISTRY"
      repository:
        registry: "quay.io"
    targets:
    - type: secret
      name: cache-secret
      namespace: vsphere-infra-helpers
      key: "the-key"
      action: "UPDATE_FIELD"
      sourceFieldIndex: 0
    - type: secret
      name: cache-secret
      namespace: vsphere-infra-helpers
      key: "the-key2"
      action: "UPDATE_FIELD"
      sourceFieldIndex: 1
    - type: deployment
      name: machine
      #namespace: vsphere-infra-helpers
      namespace: openshift-machine-api
      name: machine-api-operator
      action: "REDEPLOY_PODS"
```