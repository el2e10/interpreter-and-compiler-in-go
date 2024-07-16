package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

func Start(io io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(io)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		lxr := lexer.New(line)

		for tok := lxr.NextToken(); tok.Type != token.EOF; tok = lxr.NextToken() {
			fmt.Fprintf(out, "%v\n", tok)
		}

	}

}
