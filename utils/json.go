package utils

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func MapJSONToStringSlice(data datatypes.JSON) []string {
	if len(data) == 0 {
		return []string{}
	}

	var result []string

	if err := json.Unmarshal(data, &result); err != nil {
		return []string{}
	}

	return result
}
