package repl

import (
	"bufio"
	"fmt"
	"github.com/iamBharatManral/atom.git/cmd/internal/env"
	"github.com/iamBharatManral/atom.git/cmd/internal/interpreter"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/parser"
	"github.com/iamBharatManral/atom.git/cmd/internal/util"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

const PROMPT = "λ> "

var inputCh = make(chan string)

func Start() {
	util.Banner()
	message()
	env := env.New(nil)
	for {
		input := userInput()
		if input == "" {
			continue
		}
		lexer := lexer.New([]rune(input))
		parser := parser.New(lexer)
		program := parser.Parse()
		if len(parser.Errors) > 0 {
			for _, err := range parser.Errors {
				fmt.Println(err)
			}
			continue
		}
		for _, stmt := range program.Body {
			result := interpreter.Eval(stmt, env)
			if result.Type == "error" {
				fmt.Println(result.Value)
				continue
			} else if result.Type == "" || result.Value == "()" {
				fmt.Println()
				continue
			}
			fmt.Printf("%v\n", result.Value)

		}
	}
}

func userInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(PROMPT)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	input := scanner.Text()
	if len(input) == 0 {
		return ""
	}
	if input == ":q" || input == ":quit" {
		os.Exit(0)
	}
	if input == "clear" {
		clearTerminal()
		return ""
	}
	if input[len(input)-1] != ';' {
		fmt.Println("error: missing semicolon at the end!")
		return ""
	}
	return input

}
func message() {
	var username string
	currentUser, err := user.Current()
	if err != nil {
		username = "STRANGER"
	} else {
		username = currentUser.Username
	}
	fmt.Printf("🪐 Welcome %s! to the beginning of the 'LANGUAGE UNIVERSE' 🪐✨\n\n", username)
	fmt.Printf("To disappear from this universe, type ':q' or ':quit' 🚀\n\n")
}

func clearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
