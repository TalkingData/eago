package utils

// MergeMapStringInterface 合并多个MapStringInterfaceø
func MergeMapStringInterface(mObj ...map[string]interface{}) map[string]interface{} {
	newObj := make(map[string]interface{})
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}
