package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/vita-dounai/Firework/evaluator"
	"github.com/vita-dounai/Firework/lexer"
	"github.com/vita-dounai/Firework/object"
	"github.com/vita-dounai/Firework/parser"
)

const PROMPT = ">> "
const CONTINUE_PROMPT = ".."

func checkInputNotEnd(p *parser.Parser) bool {
	if len(p.Errors()) == 1 {
		if p.Errors()[0] == parser.UNEXPECTED_EOF {
			return true
		}
	}
	return false
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	p := parser.NewParser()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if strings.HasPrefix(line, ".") {
			command := line[1:]

			switch command {
			case "exit":
				return
			default:
				io.WriteString(out, fmt.Sprintf("Unknown command: %s\n", command))
				continue
			}
		}

		l := lexer.NewLexer(line)
		p.Init(l)
		program := p.ParseProgram()

		for checkInputNotEnd(p) {
			ident := p.Ident()
			fmt.Printf(CONTINUE_PROMPT + strings.Repeat(".", ident*2) + " ")

			scanned := scanner.Scan()
			if !scanned {
				return
			}

			restLine := scanner.Text()
			line += restLine

			l := lexer.NewLexer(line)
			p.Init(l)
			program = p.ParseProgram()
		}

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []parser.ParseError) {
	for _, err := range errors {
		io.WriteString(out, err.Type()+": "+err.Info()+"\n")
	}
}
