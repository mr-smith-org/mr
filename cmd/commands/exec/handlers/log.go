package execHandlers

import (
	"fmt"

	"github.com/arthurbcp/kuma/v2/internal/helpers"
	"github.com/arthurbcp/kuma/v2/pkg/functions"
	"github.com/arthurbcp/kuma/v2/pkg/style"
)

func HandleLog(log string, vars map[string]interface{}) error {
	var err error

	log, err = helpers.ReplaceVars(log, vars, functions.GetFuncMap())
	if err != nil {
		return fmt.Errorf("parsing log error: %s", err.Error())
	}

	style.LogPrint(log)
	return nil
}
