package helpers

import (
	"encoding/json"
	"path/filepath"

	"github.com/mr-smith-org/mr/pkg/filesystem"
	"github.com/mr-smith-org/mr/pkg/functions"
	"gopkg.in/yaml.v3"
)

func UnmarshalFileAndReplaceVars(fileName string, vars map[string]interface{}, fs filesystem.FileSystemInterface) (map[string]interface{}, error) {
	// Read the content of the OpenAPI file.
	fileContent, err := fs.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	fileContent, err = ReplaceVars(fileContent, vars, functions.GetFuncMap())
	if err != nil {
		return nil, err
	}
	// Unmarshal the JSON or YAML content into a generic map.
	fileData, err := UnmarshalByExt(fileName, []byte(fileContent))
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

func UnmarshalFile(fileName string, fs filesystem.FileSystemInterface) (map[string]interface{}, error) {
	// Read the content of the OpenAPI file.
	fileContent, err := fs.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON or YAML content into a generic map.
	fileData, err := UnmarshalByExt(fileName, []byte(fileContent))
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

func UnmarshalByExt(file string, configData []byte) (map[string]interface{}, error) {
	// Determine the file type based on its extension and unmarshal accordingly.
	switch filepath.Ext(file) {
	case ".yaml", ".yml":
		data, err := UnmarshalYaml(configData)
		if err != nil {
			return nil, err
		}
		return data, nil
	case ".json":
		data, err := UnmarshalJson(configData)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		res := make(map[string]interface{})
		res["content"] = string(configData)
		return res, nil
	}
}

// UnmarshalJson parses JSON configuration data into BuilderData.
//
// Parameters:
//   - configData: A byte slice containing JSON-formatted configuration data.
//
// Returns:
//
//	A pointer to BuilderData and an error if unmarshaling fails.
func UnmarshalJson(configData []byte) (map[string]interface{}, error) {
	fileData := make(map[string]interface{})
	err := json.Unmarshal(configData, &fileData)
	if err != nil {
		return fileData, err
	}
	return fileData, nil
}

// UnmarshalYaml parses YAML configuration data into BuilderData.
//
// Parameters:
//   - configData: A byte slice containing YAML-formatted configuration data.
//
// Returns:
//
//	A pointer to BuilderData and an error if unmarshaling fails.
func UnmarshalYaml(configData []byte) (map[string]interface{}, error) {
	fileData := make(map[string]interface{})
	err := yaml.Unmarshal(configData, &fileData)
	if err != nil {
		return fileData, err
	}
	return fileData, nil
}
