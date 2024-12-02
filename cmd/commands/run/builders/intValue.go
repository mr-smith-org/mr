package execBuilders

import (
	"strconv"
)

func BuildIntValue(key string, input map[string]interface{}, vars map[string]interface{}, required bool, component string) (int, error) {
	val := 0
	val, ok := input[key].(int)
	if !ok {
		if valStr, err := BuildStringValue(key, input, vars, required, component); err == nil {
			if valStr == "" {
				return 0, nil
			}
			val, err = strconv.Atoi(valStr)
			if err != nil {
				return 0, err
			}
		}
	}
	return val, nil
}
