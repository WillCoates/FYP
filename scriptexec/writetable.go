package main

import (
	"errors"

	lua "github.com/Shopify/go-lua"
)

func WriteTable(table map[string]interface{}, l *lua.State) error {
	// https://golang.org/pkg/encoding/json/#Unmarshal datatypes
	l.NewTable()

	for key, value := range table {
		l.PushString(key)
		switch value.(type) {
		case bool:
			l.PushBoolean(value.(bool))
		case float64:
			l.PushNumber(value.(float64))
		case string:
			l.PushString(value.(string))
		case []interface{}:
			err := WriteArray(value.([]interface{}), l)
			if err != nil {
				return err
			}
		case map[string]interface{}:
			WriteTable(value.(map[string]interface{}), l)
		default:
			return errors.New("Invalid type")
		}
		l.SetTable(-3)
	}

	return nil
}

func WriteArray(array []interface{}, l *lua.State) error {
	l.NewTable()
	for key, value := range array {
		l.PushNumber(float64(key))
		switch value.(type) {
		case bool:
			l.PushBoolean(value.(bool))
		case float64:
			l.PushNumber(value.(float64))
		case string:
			l.PushString(value.(string))
		case []interface{}:
			err := WriteArray(value.([]interface{}), l)
			if err != nil {
				return err
			}
		case map[string]interface{}:
			WriteTable(value.(map[string]interface{}), l)
		default:
			return errors.New("Invalid type")
		}
		l.SetTable(-3)
	}

	return nil
}
