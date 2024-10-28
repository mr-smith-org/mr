package handlers

import (
	"fmt"

	"github.com/arthurbcp/kuma/v2/internal/helpers"
	"github.com/arthurbcp/kuma/v2/pkg/functions"
	"github.com/arthurbcp/kuma/v2/pkg/style"
)

type LogHandler struct {
}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) Handle(data any, vars map[string]any) error {
	return handleLog(data.(string), vars)
}

func handleLog(log string, vars map[string]interface{}) error {
	var err error

	log, err = helpers.ReplaceVars(log, vars, functions.GetFuncMap())
	if err != nil {
		return fmt.Errorf("parsing log error: %s", err.Error())
	}

	style.LogPrint(log)
	return nil
}
