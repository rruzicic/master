package object

type Environment struct {
	store     map[string]Object
	enclosing *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, enclosing: nil}
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.enclosing != nil {
		obj, ok = e.enclosing.Get(name)
	}
	return obj, ok
}
