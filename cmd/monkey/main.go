package main

import (
	"fmt"
	"github.com/ssttuu/monkey/pkg/repl"
	"log"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", u.Username)
	fmt.Printf("Feel free to type in commands\n")

	repl.Start(os.Stdin, os.Stdout)
}
