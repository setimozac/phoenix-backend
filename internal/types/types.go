package types

type EnvManager struct {
	ID int `json:"id"`
	Enabled bool `json:"enabled"`
	UIEnabled bool `json:"uiEnabled"`
	MinReplica int32 `json:"minReplica"`
	Name string `json:"name"`
	LastUpdate int64 `json:"lastUpdate,omitempty"`
	*Metadata `json:"metadata"`
	Events []string `json:"events,omitempty"`
}

type Metadata struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
}

type Event struct {
	Name string `json:"name"`
	ServiceName string `json:"service_name,omitempty"`
}