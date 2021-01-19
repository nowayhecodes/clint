package main

import (
	"clint/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s, you're in Clint REPL!\n", user.Username)
	fmt.Printf("Type some commands...\n")
	repl.Start(os.Stdin, os.Stdout)
}
