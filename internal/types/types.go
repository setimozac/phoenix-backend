package types

type Service struct {
	ID int `json:"id"`
	Enable bool `json:"enable"`
	MinReplicas int32 `json:"min_replicas"`
	Name string `json:"name"`
}