package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	env := make(map[string]Object)
	return &Environment{store: env}
}

func ExtendEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) setIfExist(name string, value Object) bool {
	if _, ok := e.store[name]; ok {
		e.store[name] = value
		return true
	}

	if e.outer != nil {
		return e.outer.setIfExist(name, value)
	}

	return false
}

func (e *Environment) Set(name string, value Object) Object {
	if !e.setIfExist(name, value) {
		e.store[name] = value
	}
	return value
}
