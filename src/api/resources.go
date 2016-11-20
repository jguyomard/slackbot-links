package api

import (
	"encoding/json"
	"reflect"
)

type resource struct {
	data         interface{}
	meta         interface{}
	errorCode    int
	errorMessage string
}

type transformerFunc func(interface{}) interface{}

// Collection Resource
func collection(items interface{}, transformer transformerFunc) *resource {
	res := new(resource)

	// iterm must be a slice-type
	if reflect.TypeOf(items).Kind() != reflect.Slice {
		return res
	}

	// Transform all items
	itemsSlice := reflect.ValueOf(items)
	itemsTransformed := []interface{}{}
	for i := 0; i < itemsSlice.Len(); i++ {
		itemsTransformed = append(itemsTransformed, transformer(itemsSlice.Index(i).Interface()))
	}

	res.data = itemsTransformed
	return res
}

// Item Resource
func item(item interface{}, transformer transformerFunc) *resource {
	res := new(resource)
	res.data = transformer(item)
	return res
}

// Error Resource
func errors(errorMessage string, errorCode int) *resource {
	res := new(resource)
	res.errorMessage = errorMessage
	res.errorCode = errorCode
	return res
}

func (r *resource) SetMeta(meta interface{}) *resource {
	r.meta = meta
	return r
}

func (r *resource) ToArray() map[string]interface{} {
	if r.errorMessage != "" {
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    r.errorCode,
				"message": r.errorMessage,
			},
		}
	}
	return dataArraySerializer(r.data, r.meta)
}

func (r *resource) ToJSON() string {
	resJSON, _ := json.Marshal(r.ToArray())
	return string(resJSON)
}
