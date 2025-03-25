package object

type Environment struct {
	store     map[string]Object
	enclosing *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, enclosing: nil}
}
