package handlers

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mr-smith-org/mr/internal/domain"
	"github.com/mr-smith-org/mr/internal/helpers"
	"github.com/mr-smith-org/mr/pkg/functions"
	"github.com/mr-smith-org/mr/pkg/style"
	"github.com/spf13/afero"
)

// BuilderHandler manages the building process of the project structure.
// It interacts with the domain.Builder to create directories and files based on templates.
type BuilderHandler struct {
	// builder is the domain Builder responsible for providing structure and template data.
	builder *domain.Builder
}

// NewBuilderHandler creates and returns a new BuilderHandler instance.
//
// Parameters:
//   - builder: A pointer to the domain.Builder instance.
//
// Returns:
//
//	A pointer to a BuilderHandler.
func NewBuilderHandler(builder *domain.Builder) *BuilderHandler {
	return &BuilderHandler{
		builder: builder,
	}
}

// Build initiates the building process by applying templates and creating
// directories and files as defined in the Builder's data.
//
// Returns:
//
//	An error if the build process fails, otherwise nil.
func (h *BuilderHandler) Build() error {
	style.LogPrint("applying templates...")
	// Start recursive creation of directories and files from the root.
	err := h.createDirAndFilesRecursive("", h.builder.Data.Structure, h.builder.Config.ProjectPath)
	if err != nil {
		return err
	}

	fmt.Println()

	return nil
}

// createDirAndFilesRecursive recursively creates directories and files based on the provided structure.
//
// Parameters:
//   - key: The current directory or file name.
//   - node: The nested structure (directories or file definitions).
//   - basePath: The accumulated file system path from previous recursion levels.
//
// Returns:
//
//	An error if directory or file creation fails, otherwise nil.
func (h *BuilderHandler) createDirAndFilesRecursive(key string, node interface{}, basePath string) error {
	// Construct the current path by joining the base path with the current key.
	currentPath := filepath.Join(basePath, key)

	// Create the directory if it does not exist.
	err := h.builder.Fs.CreateDirectoryIfNotExists(currentPath)
	if err != nil {
		return err
	}

	switch children := node.(type) {
	case map[string]interface{}:
		// Iterate through the map to handle subdirectories and files.
		for childKey, childValue := range children {
			if len(strings.Split(childKey, ".")) > 1 {
				err := h.createFileAndApplyTemplate(currentPath, childKey, childValue.(map[string]interface{}))
				if err != nil {
					style.CrossMarkPrint(filepath.Join(currentPath, childKey))
					return err
				}
				style.CheckMarkPrint(filepath.Join(currentPath, childKey))
				continue
			}

			childKey, err := helpers.ReplaceVars(childKey, childValue, functions.GetFuncMap())
			if err != nil {
				return err
			}

			// Recursively create subdirectories and files.
			err = h.createDirAndFilesRecursive(childKey, childValue, currentPath)
			if err != nil {
				return err
			}
		}
	default:
	}

	return nil
}

// createFileAndApplyTemplate creates a file and applies the specified template to it.
//
// Parameters:
//   - currentPath: The directory path where the file will be created.
//   - fileName: The name of the file to be created.
//   - data: A map containing template data and metadata.
//
// Returns:
//
//	An error if file creation or template application fails, otherwise nil.
func (h *BuilderHandler) createFileAndApplyTemplate(currentPath string, fileName string, data map[string]interface{}) error {
	filePath := filepath.Join(currentPath, fileName)

	file, err := h.builder.Fs.CreateFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	t, err := h.getTemplate(data)
	if err != nil {
		return err
	}
	data = map[string]interface{}{
		"data":   data["data"],
		"global": h.builder.Data.Global,
	}
	return t.Execute(file, data)
}

// getTemplate retrieves and parses the template files based on the provided data.
//
// Parameters:
//   - data: A map containing template metadata, including the template name and any includes.
//
// Returns:
//
//	A pointer to the parsed template.Template and an error if parsing fails.
func (h *BuilderHandler) getTemplate(data map[string]interface{}) (*template.Template, error) {
	templateName, ok := data["template"].(string)
	if !ok || templateName == "" {
		return nil, fmt.Errorf("template is required")
	}

	templateFile := filepath.Join(h.builder.Config.TemplatesPath, templateName)

	allTemplates := []string{templateFile}

	if includes, ok := data["includes"].([]interface{}); ok {
		for _, include := range includes {
			includeStr, ok := include.(string)
			if !ok {
				return nil, fmt.Errorf("invalid include type: %v", include)
			}
			allTemplates = append(allTemplates, filepath.Join(h.builder.Config.TemplatesPath, includeStr))
		}
	}

	tmpl := template.New(templateName).Funcs(functions.GetFuncMap())

	for _, tmplFile := range allTemplates {
		content, err := afero.ReadFile(h.builder.Fs.GetAferoFs(), tmplFile)
		if err != nil {
			return nil, fmt.Errorf("error reading template file %s: %w", tmplFile, err)
		}

		_, err = tmpl.Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("error parsing template file %s: %w", tmplFile, err)
		}
	}

	return tmpl, nil
}
