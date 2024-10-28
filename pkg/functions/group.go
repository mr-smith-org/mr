package functions

func GroupByKey(data []interface{}, key string) map[string]interface{} {
	groupedData := make(map[string]interface{})
	for _, item := range data {
		itemMap := item.(map[string]interface{})
		keyValue := itemMap[key].(string)
		if _, ok := groupedData[keyValue]; !ok {
			groupedData[keyValue] = make([]interface{}, 0)
		}
		groupedData[keyValue] = append(groupedData[keyValue].([]interface{}), item)
	}
	return groupedData
}
