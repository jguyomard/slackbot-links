package api

import "encoding/json"

// Serializers
func dataArraySerializer(res interface{}) string {
	output := map[string]interface{}{
		"data": res,
	}
	resJSON, _ := json.Marshal(output)
	return string(resJSON)
}
