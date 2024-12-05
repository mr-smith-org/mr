package domain

type ModulePipeline struct {
	Description string `json:"description"`
	File        string `json:"file"`
}

type Module struct {
	Key         string                    `json:"key"`
	Description string                    `json:"description"`
	Version     string                    `json:"version"`
	Pipelines   map[string]ModulePipeline `json:"pipelines"`
}

func NewModule(moduleKey string, module map[string]interface{}, pipelines map[string]Pipeline) Module {
	pipelinesMap := map[string]ModulePipeline{}
	for key, pipeline := range pipelines {
		pipelinesMap[key] = ModulePipeline{
			Description: pipeline.Description,
			File:        pipeline.File,
		}
	}
	return Module{
		Key:         moduleKey,
		Description: module["description"].(string),
		Version:     module["version"].(string),
		Pipelines:   pipelinesMap,
	}
}
