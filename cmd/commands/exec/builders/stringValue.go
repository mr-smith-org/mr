package execBuilders

import (
	"fmt"

	"github.com/mr-smith/mr/internal/helpers"
	"github.com/mr-smith/mr/pkg/functions"
)

func BuildStringValue(key string, input map[string]interface{}, vars map[string]interface{}, required bool, component string) (string, error) {
	var err error
	val, ok := input[key].(string)
	if !ok {
		if required {
			return "", fmt.Errorf("%s is required for %s", key, component)
		}
		return "", nil
	}
	val, err = helpers.ReplaceVars(val, vars, functions.GetFuncMap())
	if err != nil {
		return "", err
	}
	return val, nil
}
