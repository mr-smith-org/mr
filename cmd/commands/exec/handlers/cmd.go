package handlers

import (
	"fmt"
	"strings"

	"github.com/mr-smith/mr/cmd/shared"
	"github.com/mr-smith/mr/internal/helpers"
	"github.com/mr-smith/mr/pkg/functions"
	"github.com/mr-smith/mr/pkg/style"
)

type CmdHandler struct {
}

func NewCmdHandler() *CmdHandler {
	return &CmdHandler{}
}

func (h *CmdHandler) Handle(data any, vars map[string]any) error {
	return handleCommand(data.(string), vars)
}

func handleCommand(cmdStr string, vars map[string]interface{}) error {
	var err error

	cmdStr, err = helpers.ReplaceVars(cmdStr, vars, functions.GetFuncMap())
	if err != nil {
		return fmt.Errorf("parsing command error: %s", err.Error())
	}

	style.LogPrint(fmt.Sprintf("running: %s", cmdStr))

	cmdArgs := strings.Split(cmdStr, " ")
	if err := shared.RunCommand(cmdArgs[0], cmdArgs[1:]...); err != nil {
		return fmt.Errorf("command error: %s", err.Error())
	}
	return nil
}
