package api

import "time"

func dateOrNil(val interface{}) interface{} {
	date, isDate := val.(*time.Time)
	if !isDate || date == nil {
		return nil
	}
	return date.Format("2006-01-02 15:04:05")
}

func stringOrNil(val interface{}) interface{} {
	if val == "" {
		return nil
	}
	return val
}
