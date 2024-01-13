package env

import (
	"errors"

	"github.com/nikgalushko/gan-ilox/internal"
)

var ErrUndefinedVariable = errors.New("undefined variable")

type Environment struct {
	parent    *Environment
	variables map[string]internal.Literal
}

func New() *Environment {
	return &Environment{variables: make(map[string]internal.Literal)}
}

func NewWithParent(parent *Environment) *Environment {
	return &Environment{
		variables: make(map[string]internal.Literal),
		parent:    parent,
	}
}

func (e *Environment) Get(name string) (internal.Literal, error) {
	v, ok := e.variables[name]
	if ok {
		return v, nil
	}
	if e.parent == nil {
		return internal.LiteralNil, ErrUndefinedVariable
	}

	return e.parent.Get(name)
}

func (e *Environment) Define(name string, value internal.Literal) {
	e.variables[name] = value
}

func (e *Environment) Assign(name string, value internal.Literal) error {
	if _, ok := e.variables[name]; ok {
		e.variables[name] = value
		return nil
	}

	if e.parent == nil {
		return ErrUndefinedVariable
	}

	return e.parent.Assign(name, value)
}
