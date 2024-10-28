package functions

import "strings"

func GetParamsByType(params []interface{}, paramType string) []interface{} {
	filteredParams := make([]interface{}, 0)
	for _, param := range params {
		if paramMap, ok := param.(map[string]interface{}); ok {
			if paramTypeStr, ok := paramMap["in"].(string); ok {
				if paramTypeStr == paramType {
					filteredParams = append(filteredParams, param)
				}
			}
		}
	}
	return filteredParams
}

func GetPathsByTag(paths map[string]interface{}, tag string) map[string]interface{} {
	filteredPaths := make(map[string]interface{})
	for path, pathItem := range paths {
		if pathMap, ok := pathItem.(map[string]interface{}); ok {
			for _, operation := range pathMap {
				if operationMap, ok := operation.(map[string]interface{}); ok {
					if pathTags, ok := operationMap["tags"].([]interface{}); ok {
						for _, tagItem := range pathTags {
							if tagStr, ok := tagItem.(string); ok {
								if tagStr == tag {
									filteredPaths[path] = pathItem
								}
							}
						}
					}
				}
			}
		}
	}
	return filteredPaths
}

func GetRefFrom(object map[string]interface{}) string {
	ref, ok := object["$ref"].(string)
	if !ok {
		return ""
	}
	const refPrefix = "#/definitions/"
	if strings.HasPrefix(ref, refPrefix) {
		return strings.TrimPrefix(ref, refPrefix)
	}
	return ""
}
