package api

import (
	"testing"
	"time"
)

func TestDateOrNil_nil(t *testing.T) {
	date := interface{}(nil)
	if dateOrNil(date) != nil {
		t.Fatal("TestDateOrNil_nil: date isn't nil")
	}
}

func TestDateOrNil_date(t *testing.T) {
	date := time.Now()
	if _, ok := dateOrNil(&date).(string); !ok {
		t.Fatal("TestDateOrNil_date: date isn't a string")
	}
}

func TestStringOrNil_nil(t *testing.T) {
	str := interface{}(nil)
	if stringOrNil(str) != nil {
		t.Fatal("TestStringOrNil_nil: str isn't nil")
	}
}

func TestStringOrNil_string(t *testing.T) {
	str := "test"
	if _, ok := stringOrNil(str).(string); !ok {
		t.Fatal("TestStringOrNil_string: str isn't a string")
	}
}
