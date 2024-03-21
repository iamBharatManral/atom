package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/iamBharatManral/atom.git/cmd/internal/env"
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
	env := env.New()
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
		program := parser.Parse()
		if len(program.Body) > 0 {
			result := interpreter.Eval(program.Body[0], env)
			if result.Type == "error" {
				fmt.Println(result.Value)
				continue
			} else if result.Type == "" {
				continue
			}
			fmt.Println(result.Value)
		}
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
