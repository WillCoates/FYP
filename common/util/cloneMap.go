package util

func CloneMapStringIface(toClone map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	if toClone != nil {
		for key, value := range toClone {
			newMap[key] = value
		}
	}
	return newMap
}
