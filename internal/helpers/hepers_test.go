package helpers

import (
	"reflect"
	"testing"

	"github.com/arthurbcp/kuma/v2/pkg/filesystem"
	"github.com/arthurbcp/kuma/v2/pkg/functions"
	"github.com/spf13/afero"
)

func TestUnmarshalFile(t *testing.T) {
	// Create a temporary JSON file
	aferoFS := afero.NewMemMapFs()

	jsonContent := `{"key": "value"}`
	aferoFS.Create("test.json")
	afero.WriteFile(aferoFS, "test.json", []byte(jsonContent), 0644)

	// Create a temporary YAML file
	yamlContent := "key: value"
	aferoFS.Create("test.yaml")
	afero.WriteFile(aferoFS, "test.yaml", []byte(yamlContent), 0644)

	fs := filesystem.NewFileSystem(aferoFS)

	tests := []struct {
		name     string
		fileName string
		want     map[string]interface{}
		wantErr  bool
	}{
		{"JSON file", "test.json", map[string]interface{}{"key": "value"}, false},
		{"YAML file", "test.yaml", map[string]interface{}{"key": "value"}, false},
		{"Non-existent file", "nonexistent.json", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalFile(tt.fileName, fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringContains(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		e    string
		want bool
	}{
		{"Contains", []string{"a", "b", "c"}, "b", true},
		{"Does not contain", []string{"a", "b", "c"}, "d", false},
		{"Empty slice", []string{}, "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringContains(tt.s, tt.e); got != tt.want {
				t.Errorf("StringContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceContains(t *testing.T) {
	tests := []struct {
		name string
		s    []interface{}
		e    string
		want bool
	}{
		{"Contains", []interface{}{"a", "b", "c"}, "b", true},
		{"Does not contain", []interface{}{"a", "b", "c"}, "d", false},
		{"Empty slice", []interface{}{}, "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceContains(tt.s, tt.e); got != tt.want {
				t.Errorf("InterfaceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToYaml(t *testing.T) {
	data := map[string]string{"key": "value"}
	result := functions.ToYaml(data)
	expected := []string{"key: value", ""}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToYaml() = %v, want %v", result, expected)
	}
}

func TestGetRefFrom(t *testing.T) {
	tests := []struct {
		name   string
		object map[string]interface{}
		want   string
	}{
		{"Valid ref", map[string]interface{}{"$ref": "#/definitions/Example"}, "Example"},
		{"No ref", map[string]interface{}{"key": "value"}, ""},
		{"Invalid ref", map[string]interface{}{"$ref": "Invalid"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := functions.GetRefFrom(tt.object); got != tt.want {
				t.Errorf("GetRefFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrettyJson(t *testing.T) {
	input := `{"key":"value"}`
	expected := "{\n\t\"key\": \"value\"\n}"

	result := PrettyJson(input)
	if result != expected {
		t.Errorf("PrettyJson() = %v, want %v", result, expected)
	}
}

func TestPrettyMarshal(t *testing.T) {
	input := map[string]string{"key": "value"}
	expected := "{\n\t\"key\": \"value\"\n}"

	result, err := PrettyMarshal(input)
	if err != nil {
		t.Errorf("PrettyMarshal() error = %v", err)
	}
	if result != expected {
		t.Errorf("PrettyMarshal() = %v, want %v", result, expected)
	}
}

func TestReplaceVars(t *testing.T) {
	template := "Hello, {{.Name}}!"
	vars := map[string]string{"Name": "World"}
	expected := "Hello, World!"

	result, err := ReplaceVars(template, vars, nil)
	if err != nil {
		t.Errorf("ReplaceVars() error = %v", err)
	}
	if result != expected {
		t.Errorf("ReplaceVars() = %v, want %v", result, expected)
	}
}
