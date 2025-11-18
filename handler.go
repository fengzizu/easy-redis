package main

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

var SETMap = map[string]string{}
var SETLock = sync.RWMutex{}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERROR: wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETLock.Lock()
	SETMap[key] = value
	SETLock.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERROR: wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETLock.RLock()
	value, ok := SETMap[key]
	SETLock.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var HSETMap = map[string]map[string]string{}
var HSETLock = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERROR: wrong number of arguments for 'hset' command"}
	}

	mapKey := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETLock.Lock()
	if _, ok := HSETMap[mapKey]; !ok {
		HSETMap[mapKey] = map[string]string{}
	}
	HSETMap[mapKey][key] = value
	HSETLock.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERROR: wrong number of arguments for 'hget' command"}
	}

	mapKey := args[0].bulk
	key := args[1].bulk
	HSETLock.RLock()
	value, ok := HSETMap[mapKey][key]
	HSETLock.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERROR: wrong number of arguments for 'hgetall' command"}
	}

	mapKey := args[0].bulk
	var result []Value

	HSETLock.RLock()
	value, ok := HSETMap[mapKey]
	if !ok {
		return Value{typ: "null"}
	}

	for key, item := range value {
		result = append(result, Value{typ: "bulk", bulk: fmt.Sprintf("%v: %v", key, item)})
	}

	HSETLock.RUnlock()

	return Value{typ: "array", array: result}
}
