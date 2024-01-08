package env

import (
	"errors"

	"github.com/nikgalushko/gan-ilox/token"
)

var ErrUndefinedVariable = errors.New("undefined variable")

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
	if e.parent == nil {
		return token.LiteralNil, ErrUndefinedVariable
	}

	return e.parent.Get(name)
}

func (e *Environment) Define(name string, value token.Literal) {
	e.variables[name] = value
}

func (e *Environment) Assign(name string, value token.Literal) error {
	if _, ok := e.variables[name]; ok {
		e.variables[name] = value
		return nil
	}

	if e.parent == nil {
		return ErrUndefinedVariable
	}

	return e.parent.Assign(name, value)
}
