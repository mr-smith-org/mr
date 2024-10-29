package domain

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/kuma-framework/kuma/v2/internal/helpers"
	"github.com/kuma-framework/kuma/v2/pkg/filesystem"
	"github.com/kuma-framework/kuma/v2/pkg/functions"
	"github.com/kuma-framework/kuma/v2/pkg/style"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// BuilderData encapsulates the structure and templates data parsed from configuration files.
type BuilderData struct {
	// Structure defines the directory and file hierarchy to be created.
	Structure map[string]interface{}

	// Templates defines the templates to be applied to the generated files.
	Templates map[string]interface{}

	// Global defines the global variables to be used in all the templates.
	Global map[string]interface{}
}

// Builder is responsible for managing the configuration and data required to build the project structure.
type Builder struct {
	// Config holds the configuration paths for the project and templates.
	Config *Config

	// Data holds the parsed structure and templates data.
	Data *BuilderData

	// ParsedData holds the parsed file content.
	ParsedData string

	// Fs is the file system service used to interact with the file system.
	Fs filesystem.FileSystemInterface
}

// NewBuilder initializes a new Builder instance.
//
// Parameters:
//   - file: The path to the configuration file (JSON or YAML).
//   - vars: A map of variables to replace placeholders in the configuration file.
//   - config: A pointer to the Config struct containing project and templates paths.
//
// Returns:
//
//	A pointer to a Builder instance if successful, or an error if initialization fails.
func NewBuilder(fs filesystem.FileSystemInterface, config *Config) (*Builder, error) {
	builder := Builder{}
	builder.Fs = fs
	err := builder.setConfig(config)
	if err != nil {
		return nil, err
	}
	return &builder, nil
}

// SetBuilderDataFromFile parses the configuration file and populates the BuilderData.
//
// Parameters:
//   - file: The path to the configuration file.
//   - vars: A map of variables for placeholder replacement in the configuration.
//
// Returns:
//
//	An error if parsing fails, otherwise nil.
func (b *Builder) SetBuilderDataFromFile(file string, vars map[string]interface{}) error {
	style.LogPrint("parsing config...")

	configData, err := b.Fs.ReadFile(file)
	if err != nil {
		return err
	}

	configData, err = helpers.ReplaceVars(configData, vars, functions.GetFuncMap())
	b.ParsedData = string(configData)
	style.DebugPrint("Config file", b.ParsedData)
	if err != nil {
		return err
	}

	switch filepath.Ext(file) {
	case ".yaml", ".yml":
		data, err := unmarshalYamlConfig([]byte(configData))
		if err != nil {
			return err
		}
		b.Data = data

	case ".json":
		data, err := unmarshalJsonConfig([]byte(configData))
		if err != nil {
			return err
		}
		b.Data = data
	default:
		return fmt.Errorf("invalid file extension: %s", file)
	}
	return nil
}

// SetConfig assigns the provided Config to the Builder.
//
// Parameters:
//   - config: A pointer to the Config struct.
//
// Returns:
//
//	An error if setting the configuration fails, otherwise nil.
func (b *Builder) setConfig(config *Config) error {
	b.Config = config
	return nil
}

// UnmarshalJsonConfig parses JSON configuration data into BuilderData.
//
// Parameters:
//   - configData: A byte slice containing JSON-formatted configuration data.
//
// Returns:
//
//	A pointer to BuilderData and an error if unmarshaling fails.
func unmarshalJsonConfig(configData []byte) (*BuilderData, error) {
	config := BuilderData{}
	err := json.Unmarshal(configData, &config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}

// UnmarshalYamlConfig parses YAML configuration data into BuilderData.
//
// Parameters:
//   - configData: A byte slice containing YAML-formatted configuration data.
//
// Returns:
//
//	A pointer to BuilderData and an error if unmarshaling fails.
func unmarshalYamlConfig(configData []byte) (*BuilderData, error) {
	config := BuilderData{}
	c := map[interface{}]interface{}{}
	err := yaml.Unmarshal(configData, &c)
	if err != nil {
		return &config, err
	}
	err = mapstructure.Decode(c, &config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}
