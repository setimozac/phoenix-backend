package types

type EnvManager struct {
	ID int `json:"id"`
	Enabled bool `json:"enabled"`
	UIEnabled bool `json:"ui_enabled"`
	MinReplicas int32 `json:"min_replicas"`
	Name string `json:"name"`
	LastUpdate int64 `json:"last_update"`
}