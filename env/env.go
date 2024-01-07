package env

import (
	"errors"

	"github.com/nikgalushko/gan-ilox/token"
)

type Environment struct {
	variables map[string]token.Literal
}

func New() *Environment {
	return &Environment{variables: make(map[string]token.Literal)}
}

func (e *Environment) Get(name string) (token.Literal, error) {
	v, ok := e.variables[name]
	if !ok {
		return token.LiteralNil, errors.New("undefined variable")
	}

	return v, nil
}

func (e *Environment) Set(name string, value token.Literal) {
	e.variables[name] = value
}
