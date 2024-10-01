package utils

import (
	"encoding/json"
	"time"
)

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateStr)
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}
