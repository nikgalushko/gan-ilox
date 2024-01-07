package env

import (
	"errors"

	"github.com/nikgalushko/gan-ilox/token"
)

type Environment struct {
	parent    *Environment
	variables map[string]token.Literal
}

func New() *Environment {
	return &Environment{variables: make(map[string]token.Literal)}
}

func NewWithParent(parent *Environment) *Environment {
	return &Environment{
		variables: make(map[string]token.Literal),
		parent:    parent,
	}
}

func (e *Environment) Get(name string) (token.Literal, error) {
	v, ok := e.variables[name]
	if ok {
		return v, nil
	}

	return e.parent.Get(name)
}

func (e *Environment) Define(name string, value token.Literal) {
	e.variables[name] = value
}

func (e *Environment) Assign(name string, value token.Literal) error {
	if e.Has(name) {
		e.variables[name] = value
		return nil
	}

	if e.parent == nil {
		return errors.New("undefined variable")
	}
	return e.parent.Assign(name, value)
}

func (e *Environment) Has(name string) bool {
	_, ok := e.variables[name]
	return ok
}
