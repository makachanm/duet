package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const PROMPT = ">> "

func FileExecute(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	l := New(string(file))
	p := NewParser(l)

	program := p.ParseProgram()
	memory := NewMemory()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	engine := NewExcutionEngine(program, memory)
	evaluated := engine.Run()

	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

func Repl(in io.Reader, out io.Writer) {
	fmt.Println("Duet REPL. Ctrl-C to exit.")

	scanner := bufio.NewScanner(in)
	memory := NewMemory()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := New(line)
		p := NewParser(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		engine := NewExcutionEngine(program, memory)
		evaluated := engine.Run()

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func main() {
	if len(os.Args) > 1 {
		FileExecute(os.Args[1])
		return
	}

	Repl(os.Stdin, os.Stdout)
}
