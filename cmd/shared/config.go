package shared

var (
	FilesPath       string = ".mr-smith"
	PipelinesPath   string = FilesPath + "/pipelines"
	ModulesFileName        = "mr-modules.yaml"
	ConfigFileName         = "mr-config.yaml"
	GitHubURL              = "https://github.com"

	Pipeline string
	Module   string
	Vars     = make(map[string]string)
)
