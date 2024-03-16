package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/iamBharatManral/atom.git/cmd/internal/interpreter"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/parser"
	"github.com/iamBharatManral/atom.git/cmd/internal/util"
)

const PROMPT = "λ> "

func Start() {
	util.Banner()
	message()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(PROMPT)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		input := []rune(scanner.Text())
		if string(input) == ":q" || string(input) == ":quit" {
			os.Exit(0)
		}
		lexer := lexer.New(input)
		parser := parser.New(lexer)
		parser.Parse()
		interpreter.Eval(parser.Ast)
	}
}

func message() {
	var username string
	currentUser, err := user.Current()
	if err != nil {
		username = "STRANGER"
	} else {
		username = currentUser.Username
	}
	fmt.Printf("Welcome %s to the beginning of the language universe\n\n", username)
}
