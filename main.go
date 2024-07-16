package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {

	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to the monkey language %s", user.Username)
	repl.Start(os.Stdin, os.Stdout)

}
