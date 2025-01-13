package data

type SyncSchemaDockerConfig struct {
	Registry string `yaml:"registry"`
}

type SyncItemSource struct {
	Namespace  string                  `yaml:"namespace"`
	Name       string                  `yaml:"name"`
	Type       string                  `yaml:"type"`
	Key        string                  `yaml:"key"`
	Schema     string                  `yaml:"schema"`
	Repository *SyncSchemaDockerConfig `yaml:"repository"`
}

type TriggerType string
type TargetAction string
type SyncSchema string

const (
	SYNC_TRIGGER_ON_CHANGE           = TriggerType("ON_CHANGE")
	SYNC_SCHEMA_REGISTRY             = TriggerType("REGISTRY")
	SYNC_SCHEMA_GENERIC              = TriggerType("GENERIC")
	SYNC_TARGET_ACTION_UPDATE_FIELD  = TargetAction("UPDATE_FIELD")
	SYNC_TARGET_ACTION_REDEPLOY_PODS = TargetAction("REDEPLOY_PODS")
)

type SyncItemTarget struct {
	Namespace        string `yaml:"namespace"`
	Name             string `yaml:"name"`
	Type             string `yaml:"type"`
	Key              string `yaml:"key"`
	SourceFieldIndex int    `yaml:"sourceFieldIndex"`
}

type SyncTriggerAction struct {
	Type    TriggerType      `yaml:"type"`
	Source  SyncItemSource   `yaml:"source"`
	Targets []SyncItemTarget `yaml:"targets"`
}

type SyncActions struct {
	RefreshPeriodSeconds int64               `yaml:"refreshPeriodSeconds"`
	Actions              []SyncTriggerAction `yaml:"actions"`
}

type SyncConfig struct {
	Sync   SyncActions `yaml:"sync"`
	DryRun bool
}
