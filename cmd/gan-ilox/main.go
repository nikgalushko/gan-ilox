package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nikgalushko/gan-ilox/debug"
	"github.com/nikgalushko/gan-ilox/parser"
	"github.com/nikgalushko/gan-ilox/scanner"
)

func main() {
	var err error
	args := os.Args[1:] // cut programm name

	if len(args) > 1 {
		fmt.Println("Usage: gan-ilox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		err = runFile(args[0])
	} else {
		err = runPrompt()
	}

	if err != nil {
		fmt.Println(err.Error())
	}
}

func runFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return run(string(data))
}

func runPrompt() error {
	s := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")

	for s.Scan() {
		err := run(s.Text())
		if err != nil {
			return err
		}
		fmt.Print("> ")
	}

	return s.Err()
}

func run(source string) error {
	s := scanner.NewScanner(source)
	tokens, err := s.ScanTokens()
	if err != nil {
		return err
	}

	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		return err
	}

	fmt.Println(debug.AstPrinter{E: expr})

	return nil
}
