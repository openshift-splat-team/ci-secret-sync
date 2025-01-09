# ci-secret-sync

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