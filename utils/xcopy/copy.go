package xcopy

import "encoding/json"

func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}
func AppendStrings(a, b []string) []string {
	for i := 0; i < len(b); i++ {
		a = append(a, b[i])
	}
	return a
}

func Convert(data interface{}, target interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &target)
	if err != nil {
		return err
	}

	return err
}
