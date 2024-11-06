package shared

import "github.com/mr-smith-org/mr/internal/domain"

var Modules = map[string]domain.Template{
	"mr-smith-org/typescript-rest-services": domain.NewTemplate(
		"TypeScript Rest Services (OpenAPI 2.0)",
		"Create a library TypeScript with services typed for all endpoints described in a file Open API 2.0",
		[]string{"typescript", "openapi", "rest", "library"},
	),
	"mr-smith-org/hello-world": domain.NewTemplate(
		"Hello World",
		"A simple Hello World in Go!",
		[]string{"golang", "example"},
	),
	"mr-smith-org/changelog-generator": domain.NewTemplate(
		"Changelog Generator",
		"Helper to write a good changelog to your project",
		[]string{"changelog", "helper", "markdown"},
	),
	"mr-smith-org/commit-standardizer": domain.NewTemplate(
		"Commit Standardizer",
		"Write conventional commits for your projects",
		[]string{"git", "standardizer"},
	),
}
