package api

func dataArraySerializer(data interface{}, meta interface{}) map[string]interface{} {
	output := map[string]interface{}{
		"data": data,
	}
	if meta != nil {
		output["meta"] = meta
	}
	return output
}
