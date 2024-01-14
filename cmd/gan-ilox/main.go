package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/nikgalushko/gan-ilox/debug"
	"github.com/nikgalushko/gan-ilox/env"
	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/nikgalushko/gan-ilox/interpreter"
	"github.com/nikgalushko/gan-ilox/parser"
	"github.com/nikgalushko/gan-ilox/scanner"
)

func main() {
	var err error
	args := os.Args[1:] // cut programm name
	environment := env.New()
	environment.Define("now", internal.NewLiteralNativeFunction(
		nil, func(args ...internal.Literal) (internal.Literal, error) {
			return internal.NewLiteralInt(time.Now().UnixMilli()), nil
		},
	))

	environment.Define("sleep", internal.NewLiteralNativeFunction(
		[]string{"seconds"}, func(args ...internal.Literal) (internal.Literal, error) {
			if len(args) != 1 {
				return internal.LiteralNil, errors.New("expect 1 argument got 0")
			}

			if !args[0].IsNumber() {
				return internal.LiteralNil, errors.New("expect number as argument")
			}

			time.Sleep(time.Duration(args[0].AsInt()) * time.Second)
			return internal.LiteralNil, nil
		},
	))

	if len(args) > 1 {
		fmt.Println("Usage: gan-ilox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		err = runFile(environment, args[0])
	} else {
		err = runPrompt(environment)
	}

	if err != nil {
		fmt.Println(err.Error())
	}
}

func runFile(env *env.Environment, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return run(env, string(data), false)
}

func runPrompt(env *env.Environment) error {
	s := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")

	for s.Scan() {
		err := run(env, s.Text(), true)
		if err != nil {
			return err
		}
		fmt.Print("> ")
	}

	return s.Err()
}

func run(env *env.Environment, source string, isPrompt bool) error {
	s := scanner.NewScanner(source)
	tokens, err := s.ScanTokens()
	if err != nil {
		return err
	}

	p := parser.New(tokens)
	stmts, err := p.Parse()
	if err != nil {
		return err
	}

	fmt.Println("__debug__", debug.AstPrinter{S: stmts})
	i := interpreter.New(env, stmts)
	ret, err := i.Interpret()
	if err != nil {
		fmt.Println(err.Error())
	}

	if len(ret) != 0 && isPrompt {
		for _, r := range ret {
			fmt.Println(r)
		}
	}

	return nil
}
