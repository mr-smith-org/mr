package domain

import (
	"reflect"
	"testing"

	"github.com/kuma-framework/kuma/v2/pkg/filesystem"
	"github.com/spf13/afero"
)

var aferoFs = afero.NewMemMapFs()
var mockFs = filesystem.NewFileSystem(aferoFs)

func TestNewBuilder(t *testing.T) {
	aferoFs.Create("test.yaml")
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "Successful initialization",
			config:  &Config{ProjectPath: "path"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBuilder(mockFs, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewBuilder() returned nil, want non-nil")
			}
		})
	}
}

func TestBuilder_SetBuilderData(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		content  string
		vars     map[string]interface{}
		wantErr  bool
		wantData *BuilderData
	}{
		{
			name:    "Valid YAML file",
			file:    "test.yaml",
			content: "structure:\n  dir: file\ntemplates:\n  tpl: content\nglobal:\n  var: value",
			vars:    map[string]interface{}{"key": "value"},
			wantErr: false,
			wantData: &BuilderData{
				Structure: map[string]interface{}{"dir": "file"},
				Templates: map[string]interface{}{"tpl": "content"},
				Global:    map[string]interface{}{"var": "value"},
			},
		},
		{
			name:    "Valid JSON file",
			file:    "test.json",
			content: `{"structure": {"dir": "file"}, "templates": {"tpl": "content"}, "global": {"var": "value"}}`,
			vars:    map[string]interface{}{"key": "value"},
			wantErr: false,
			wantData: &BuilderData{
				Structure: map[string]interface{}{"dir": "file"},
				Templates: map[string]interface{}{"tpl": "content"},
				Global:    map[string]interface{}{"var": "value"},
			},
		},
		{
			name:    "Invalid file extension",
			file:    "test.txt",
			content: "some content",
			vars:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "Invalid YAML content",
			file:    "invalid.yaml",
			content: "invalid: : yaml",
			vars:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "Invalid JSON content",
			file:    "invalid.json",
			content: `{"invalid": json}`,
			vars:    map[string]interface{}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new in-memory file system for each test
			memFs := afero.NewMemMapFs()
			mockFs := filesystem.NewFileSystem(memFs)

			// Write the test content to the mock file system
			err := afero.WriteFile(memFs, tt.file, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			b := &Builder{
				Fs: mockFs,
			}
			err = b.SetBuilderDataFromFile(tt.file, tt.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.SetBuilderData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(b.Data, tt.wantData) {
				t.Errorf("Builder.SetBuilderData() got = %v, want %v", b.Data, tt.wantData)
			}
		})
	}
}

func TestBuilder_setConfig(t *testing.T) {
	b := &Builder{}
	config := &Config{ProjectPath: "test"}

	err := b.setConfig(config)
	if err != nil {
		t.Errorf("Builder.setConfig() error = %v, wantErr %v", err, false)
	}

	if !reflect.DeepEqual(b.Config, config) {
		t.Errorf("Builder.setConfig() didn't set the config correctly")
	}
}

func TestUnmarshalJsonConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *BuilderData
		wantErr bool
	}{
		{
			name:  "Valid JSON",
			input: `{"structure": {"dir": "file"}, "templates": {"tpl": "content"}, "global": {"var": "value"}}`,
			want: &BuilderData{
				Structure: map[string]interface{}{"dir": "file"},
				Templates: map[string]interface{}{"tpl": "content"},
				Global:    map[string]interface{}{"var": "value"},
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"invalid": json}`,
			want:    &BuilderData{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalJsonConfig([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalJsonConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unmarshalJsonConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalYamlConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *BuilderData
		wantErr bool
	}{
		{
			name: "Valid YAML",
			input: `
structure:
  dir: file
templates:
  tpl: content
global:
  var: value
`,
			want: &BuilderData{
				Structure: map[string]interface{}{"dir": "file"},
				Templates: map[string]interface{}{"tpl": "content"},
				Global:    map[string]interface{}{"var": "value"},
			},
			wantErr: false,
		},
		{
			name:    "Invalid YAML",
			input:   "invalid: yaml: content",
			want:    &BuilderData{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalYamlConfig([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalYamlConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unmarshalYamlConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
