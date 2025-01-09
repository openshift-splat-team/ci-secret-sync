package data

type SyncItemSource struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Key       string `yaml:"key"`
}

type TriggerType string
type TargetAction string

const (
	SYNC_TRIGGER_ON_CHANGE           = TriggerType("ON_CHANGE")
	SYNC_TARGET_ACTION_UPDATE_FIELD  = TargetAction("UPDATE_FIELD")
	SYNC_TARGET_ACTION_REDEPLOY_PODS = TargetAction("REDEPLOY_PODS")
)

type SyncItemTarget struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Key       string `yaml:"key"`
}

type SyncTriggerAction struct {
	Type    TriggerType      `yaml:"type"`
	Source  SyncItemSource   `yaml:"source"`
	Targets []SyncItemTarget `yaml:"targets"`
}

type SyncActions struct {
	Actions []SyncTriggerAction `yaml:"actions"`
}

type SyncConfig struct {
	Sync SyncActions `yaml:"sync"`
}
