# ci-secret-sync

Periodically checks a secret for an update. If the secret is updated, targets are updated to reflect the update.  Optionally, pods associated with a daemonset or deployment are updated.
Actions are documented in a file in the working directory named `sync.yaml`. 

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