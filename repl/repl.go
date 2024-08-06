package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/evaluator"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

const MONKEY_FACE = `            
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			print_end_message(out)
			break
		}

		lxr := lexer.New(line)
		prsr := parser.New(lxr)

		program := prsr.ParseProgram()

		if len(prsr.Errors()) != 0 {
			print_parser_errors(out, prsr.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func print_end_message(out io.Writer) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Good Bye!!\n")
}

func print_parser_errors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "- "+msg+"\n")
	}
}
