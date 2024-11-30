package shared

var (
	FilesPath       string = ".mr-smith"
	RunsPath        string = FilesPath + "/runs"
	ModulesFileName        = "mr-modules.yaml"
	ConfigFileName         = "mr-config.yaml"
	GitHubURL              = "https://github.com"

	Run    string
	Module string
	Vars   = make(map[string]string)
)
