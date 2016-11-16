package api

import "reflect"

type transformerFunc func(interface{}) interface{}

// Collection Resource
func collection(items interface{}, transformer transformerFunc) string {

	// iterm must be a slice-type
	if reflect.TypeOf(items).Kind() != reflect.Slice {
		return ""
	}

	// Transform all items
	itemsSlice := reflect.ValueOf(items)
	itemsTransformed := []interface{}{}
	for i := 0; i < itemsSlice.Len(); i++ {
		itemsTransformed = append(itemsTransformed, transformer(itemsSlice.Index(i).Interface()))
	}

	return dataArraySerializer(itemsTransformed)
}

// Item Resource
func item(item interface{}, transformer transformerFunc) string {
	itemTransformed := transformer(item)
	return dataArraySerializer(itemTransformed)
}
