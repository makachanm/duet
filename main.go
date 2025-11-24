package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const VERSION = "0.1"
const PROMPT = "? "

func FileExecute(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	l := New(string(file))
	p := NewParser(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	engine := NewExcutionEngine(program, nil)
	evaluated := engine.Run()

	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

func Repl(in io.Reader, out io.Writer) {
	fmt.Printf("Duet version %s. Ctrl-C to exit.\n", VERSION)

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
	var version bool
	flag.BoolVar(&version, "version", false, "print Duet version")
	flag.BoolVar(&version, "v", false, "print Duet version (shorthand)")
	flag.Parse()

	if version {
		fmt.Printf("Duet version %s\n", VERSION)
		return
	}

	if flag.NArg() > 0 {
		FileExecute(flag.Arg(0))
	} else {
		Repl(os.Stdin, os.Stdout)
	}
}
