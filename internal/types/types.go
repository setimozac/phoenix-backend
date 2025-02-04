package types

type EnvManager struct {
	ID int `json:"id"`
	Enabled *bool `json:"enabled"`
	UIEnabled *bool `json:"uiEnabled"`
	MinReplica *int32 `json:"minReplica"`
	Name string `json:"name"`
	LastUpdate *int64 `json:"lastUpdate,omitempty"`
	*Metadata `json:"metadata"`
}

type Metadata struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
}