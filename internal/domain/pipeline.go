package domain

type Pipeline struct {
	Key         string        `json:"key"`
	Description string        `json:"description"`
	Steps       []interface{} `json:"steps"`
	File        string        `json:"file"`
	Visible     bool          `json:"visible"`
}

func NewPipeline(key string, description string, steps []interface{}, file string, visible bool) Pipeline {
	return Pipeline{
		Key:         key,
		Description: description,
		Steps:       steps,
		Visible:     visible,
		File:        file,
	}
}
