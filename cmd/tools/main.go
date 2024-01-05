package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("Usage: gan-tools <output dir>")
		os.Exit(64)
	}

	outputDir := args[0]

	defineAst(outputDir, "Expr", []string{
		"Binary : Left Expr[R], Operator token.Token, Right Expr[R]",
		"Grouping : Expression Expr[R]",
		"Literal : Value any",
		"Unary : Operator token.Token, Right Expr[R]",
	})
}

func defineAst(outputDir, base string, types []string) {
	path := filepath.Join(outputDir, strings.ToLower(base)+".go")
	out := bytes.NewBuffer(nil)

	fmt.Fprintln(out, "package expr")
	fmt.Fprintln(out, "")

	fmt.Fprintln(out, "import (")
	fmt.Fprintln(out, `"github.com/nikgalushko/gan-ilox/token"`)
	fmt.Fprintln(out, ")")

	fmt.Fprintf(out, "type %s[R any] interface {\n", base)
	fmt.Fprintln(out, "Accept(visitor Visitor[R]) R")
	fmt.Fprintln(out, "}")
	fmt.Fprintln(out, "")

	defineVisitor(out, base, types)

	for _, t := range types {
		tokens := strings.Split(t, ":")
		if len(tokens) != 2 {
			panic(fmt.Sprintf("invalid rule: %s", t))
		}

		structName := strings.TrimSpace(tokens[0])
		fields := strings.TrimSpace(tokens[1])
		defineType(out, structName, base, fields)
	}

	src, err := format.Source(out.Bytes())
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, _ = f.Write(src)
	_ = f.Sync()
	_ = f.Close()
}

func defineVisitor(out io.Writer, base string, types []string) {
	fmt.Fprintln(out, "type Visitor[R any] interface {")
	for _, t := range types {
		structName := strings.TrimSpace(strings.Split(t, ":")[0])
		fmt.Fprintf(out, "Visit%s%s(expr %s[R]) R\n", structName, base, structName)
	}
	fmt.Fprintf(out, "}\n\n")
}

func defineType(out io.Writer, name, base, fieldList string) {
	fmt.Fprintf(out, "type %s[R any] struct {\n", name)
	fields := strings.Split(fieldList, ",")
	for _, field := range fields {
		tokens := strings.Split(strings.TrimSpace(field), " ")
		fmt.Fprintf(out, "%s %s\n", strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1]))
	}
	fmt.Fprintf(out, "}\n")

	fmt.Fprintf(out, "func (e %s[R]) Accept(visitor Visitor[R]) R {\n", name)
	fmt.Fprintf(out, "return visitor.Visit%s%s(e)\n", name, base)
	fmt.Fprintf(out, "}\n\n")
}
