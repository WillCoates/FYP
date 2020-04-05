package main

import lua "github.com/Shopify/go-lua"

func ReadTable(pos int, l *lua.State) (map[string]interface{}, bool) {
	l.PushNil()
	result := make(map[string]interface{})
	if pos < 0 {
		pos--
	}
	for l.Next(pos) {
		// -1 is value, -2 is key
		key, ok := l.ToString(-2)
		if !ok {
			l.PushString("publish expects string keys for map, got " + l.TypeOf(-2).String())
			l.Error()
			return nil, false
		}
		var value interface{}
		switch {
		case l.IsNil(-1):
			// nop
		case l.IsBoolean(-1):
			value = l.ToBoolean(-1)
		case l.IsNumber(-1):
			value, _ = l.ToNumber(-1)
		case l.IsString(-1):
			value, _ = l.ToString(-1)
		case l.IsTable(-1):
			value, ok = ReadTable(-1, l)
			if !ok {
				return nil, false
			}
		default:
			l.PushString("publish unexpected value for map, expects bool, number, string or map, got " + l.TypeOf(-2).String())
			l.Error()
			return nil, false
		}
		result[key] = value
		l.Pop(1) // Pop value off stack, key needed for next
	}
	return result, true
}
