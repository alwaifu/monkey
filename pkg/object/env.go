package object

import "errors"

func NewEnviromentFromMap(dir map[string]interface{}) (*Environment, error) {
	store := make(map[string]Object, len(dir))
	for k, v := range dir {
		switch v := v.(type) {
		case bool:
			store[k] = Boolean(v)
		case int:
			store[k] = Integer(v)
		case int64:
			store[k] = Integer(v)
		case string:
			store[k] = String(v)
		default:
			return nil, errors.New("invalid type")
		}
	}
	return &Environment{store: store}, nil
}
func ToGoValue(obj Object) (interface{}, error) {
	switch obj := obj.(type) {
	case Boolean:
		return bool(obj), nil
	case Integer:
		return int64(obj), nil
	case String:
		return string(obj), nil
	default:
		return nil, errors.New("invalid type")
	}
}
func NewEnviroment() *Environment {
	return &Environment{store: make(map[string]Object)}
}
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnviroment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
