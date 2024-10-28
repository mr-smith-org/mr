package functions

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func ToYaml(data interface{}) []string {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		panic("Error parsing YAML template: " + err.Error())
	}
	lines := strings.Split(string(yamlData), "\n")
	return lines
}
