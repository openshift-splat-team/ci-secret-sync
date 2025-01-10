# ci-secret-sync

Periodically checks a secret for an update. If the secret is updated, targets are updated to reflect the update.  Optionally, pods associated with a daemonset are updated.
Actions are documented in a file in the working directory named `sync.yaml`. 

sample: 
```
sync:
  actions:
  - source:
      name: "pull-secret"
      namespace: "openshift-config"
      key: ".dockerconfigjson"
      type: "ON_CHANGE"
    targets:
    - type: secret
      name: cache-secret
      namespace: vsphere-infra-helpers
      key: "the-key"
      action: "UPDATE_FIELD"
    - type: daemonset
      name: cache
      #namespace: vsphere-infra-helpers
      namespace: openshift-dns
      name: dns-default
      action: "REDEPLOY_PODS"
```