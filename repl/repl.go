package repl

import (
	"bufio"
	"clint/lexer"
	"clint/parser"
	"fmt"
	"io"
)

const CLINT = `
  ___  __    __  __ _  ____    ____  ____  ____  __
 / __)(  )  (  )(  ( \(_  _)  (  _ \(  __)(  _ \(  )
( (__ / (_/\ )( /    /  )(     )   / ) _)  ) __// (_/\
 \___)\____/(__)\_)__) (__)   (__\_)(____)(__)  \____/
`

// PROMPT ...
const PROMPT = ">> "

// Start ...
func Start(in io.Reader, out io.Writer) {
	fmt.Println(CLINT)
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, error := range errors {
		io.WriteString(out, "\t"+error+"\n")
	}
}
