package shared

import "github.com/charmbracelet/huh"

type Handler interface {
	Handle(data any, vars map[string]any) error
}

type FormFieldHandler interface {
	Handle(data any, vars map[string]any) (huh.Field, string, any, error)
}
