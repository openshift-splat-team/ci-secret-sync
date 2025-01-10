package schema

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/openshift-splat-team/ci-secret-sync/data"
)

const (
	REGISTRY_CREDENTIAL_FIELD_SEPERATOR = ":"
	REGISTRY_CREDENTIALS
)

type RegistryCredentials struct {
	SchemaInterface

	Data        []byte
	Config      *data.SyncItemSource
	credentials string
}

func (r *RegistryCredentials) GetFields() ([]string, error) {

	if len(r.credentials) == 0 {
		var dockerConfig map[string]interface{}

		err := json.Unmarshal(r.Data, &dockerConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal dockerconfig json: %v", err)
		}

		if authsIface, ok := dockerConfig["auths"]; ok {
			auths := authsIface.(map[string]interface{})
			if r.Config.Repository == nil {
				return nil, fmt.Errorf("repository not defined")
			}
			fmt.Printf("repo: %v", r.Config.Repository)
			if registryIface, ok := auths[r.Config.Repository.Registry]; ok {
				registry := registryIface.(map[string]interface{})
				b64creds := registry["auth"]
				decoded, err := base64.StdEncoding.DecodeString(b64creds.(string))
				if err != nil {
					return nil, fmt.Errorf("error decoding Base64 auth: %v", err)

				}
				r.credentials = string(decoded)
			}
		} else {
			return nil, errors.New("unable to dockerconfig has no auths stanza")
		}

	}

	fmt.Printf("credentials: %s\n", r.credentials)
	field := string(r.credentials)
	fields := strings.Split(field, REGISTRY_CREDENTIAL_FIELD_SEPERATOR)

	return fields, nil
}

func (r *RegistryCredentials) GetField(idx int) (string, error) {

	fields, err := r.GetFields()
	if err != nil {
		return "", fmt.Errorf("unable to get fields: %v", err)
	}

	if idx >= len(fields) {
		return "", fmt.Errorf("field index is out of range")
	}
	return string(fields[idx]), nil
}
