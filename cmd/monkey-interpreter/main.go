package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/marcusw0/monkey-interpreter/evaluator"
	"github.com/marcusw0/monkey-interpreter/lexer"
	"github.com/marcusw0/monkey-interpreter/object"
	"github.com/marcusw0/monkey-interpreter/parser"
	"github.com/marcusw0/monkey-interpreter/repl"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Monkey programming language!\n",
			user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	if len(args) == 1 {
		os.Exit(runFile(os.Stderr, os.Stdout, args[0]))
	} else {
		fmt.Fprint(os.Stderr, "monkey can only accept 1 file\n")
		os.Exit(1)
	}
}

func runFile(errOut, out io.Writer, args string) int {
	if ok := strings.HasSuffix(args, ".mky"); !ok {
		fmt.Fprint(errOut, "monkey files end with .mky\n")
		return 1
	}

	file, err := os.ReadFile(args)
	if err != nil {
		fmt.Fprintf(errOut, "file read: %v", err)
		return 1
	}
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	l := lexer.New(string(file))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(errOut, p.Errors())
		return 1
	}

	evaluator.DefineMacros(program, macroEnv)
	expanded := evaluator.ExpandMacros(program, macroEnv)

	evaluated := evaluator.Eval(expanded, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}

	return 0
}
