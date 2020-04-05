package main

import (
	"testing"

	"github.com/Shopify/go-lua"
)

func TestEmptyTable(t *testing.T) {
	l := lua.NewState()
	l.NewTable()
	table, ok := ReadTable(-1, l)
	if !ok {
		t.Fatal("Table failed to read")
	}
	if len(table) != 0 {
		t.Fatal("Table not empty")
	}
}

func TestSimpleStringTable(t *testing.T) {
	l := lua.NewState()
	l.NewTable()
	l.PushString("a")
	l.PushString("foo")
	l.SetTable(-3)
	l.PushString("b")
	l.PushString("bar")
	l.SetTable(-3)
	table, ok := ReadTable(-1, l)
	if !ok || len(table) != 2 {
		t.Fatal("Table failed to read or wrong length")
	}
	rawValue, ok := table["a"]
	if !ok {
		t.Fatal("Failed to read key from table")
	}
	value, ok := rawValue.(string)
	if !ok || value != "foo" {
		t.Fatal("Failed to read string from table")
	}

	rawValue, ok = table["b"]
	if !ok {
		t.Fatal("Failed to read second key from table")
	}
	value, ok = rawValue.(string)
	if !ok || value != "bar" {
		t.Fatal("Failed to read second string from table")
	}
}

func TestAllTableTypes(t *testing.T) {
	l := lua.NewState()
	/*
		{
			"a": "foo",
			"b": 42,
			"c": true,
			"d": {
				"a": "bar"
			}
		}
	*/
	l.NewTable()
	l.PushString("a")
	l.PushString("foo")
	l.SetTable(-3)
	l.PushString("b")
	l.PushNumber(42)
	l.SetTable(-3)
	l.PushString("c")
	l.PushBoolean(true)
	l.SetTable(-3)
	l.PushString("d")
	l.NewTable()
	l.PushString("a")
	l.PushString("bar")
	l.SetTable(-3)
	l.SetTable(-3)

	table, ok := ReadTable(-1, l)
	if !ok || len(table) != 4 {
		t.Fatal("Failed to read table or table wrong length")
	}
	rawValue, ok := table["a"]
	if !ok {
		t.Fatal("Failed to read string key from table")
	}
	strValue, ok := rawValue.(string)
	if !ok || strValue != "foo" {
		t.Fatal("Failed to read string value from table")
	}

	rawValue, ok = table["b"]
	if !ok {
		t.Fatal("Failed to read number key from table")
	}
	numValue, ok := rawValue.(float64)
	if !ok || int64(numValue) != 42 {
		t.Fatal("Failed to read number value from table")
	}

	rawValue, ok = table["c"]
	if !ok {
		t.Fatal("Failed to read bool key from table")
	}
	boolValue, ok := rawValue.(bool)
	if !ok || !boolValue {
		t.Fatal("Failed to read bool value from table")
	}

	rawValue, ok = table["d"]
	if !ok {
		t.Fatal("Failed to read table key from table")
	}
	tabValue, ok := rawValue.(map[string]interface{})
	if !ok || len(tabValue) != 1 {
		t.Fatal("Failed to read table value from table or wrong length")
	}

	rawValue, ok = tabValue["a"]
	if !ok {
		t.Fatal("Failed to read string key from subtable")
	}
	strValue, ok = rawValue.(string)
	if !ok || strValue != "bar" {
		t.Fatal("Failed to read string value from subtable")
	}
}
