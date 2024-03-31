package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/iamBharatManral/atom.git/cmd/internal/env"
	"github.com/iamBharatManral/atom.git/cmd/internal/interpreter"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/parser"
	"github.com/iamBharatManral/atom.git/cmd/internal/util"
)

const MAIN_PROMPT = "λ> "
const REST_OF_LINE_PROMPT = "... "

func Start() {
	util.Banner()
	message()
	env := env.New()
	for {
		input := userInput()
		lexer := lexer.New(input)
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

func userInput() []rune {
	endOfStatements := make(map[string]string)
	endOfStatements["let"] = ";"
	endOfStatements["fn"] = "end;"
	var finalInput string
	scanner := bufio.NewScanner(os.Stdin)
	var endOfStatement string = endOfStatements["let"]
	var lineContinuation bool
	for {
		if !lineContinuation {
			fmt.Print(MAIN_PROMPT)
		} else {
			fmt.Print(REST_OF_LINE_PROMPT)
		}
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		input := scanner.Text()
		finalInput += input
		if strings.Contains(finalInput, "fn") {
			endOfStatement = endOfStatements["fn"]
		} else if strings.Contains(finalInput, "let") {
			endOfStatement = endOfStatements["let"]
		} else {
			endOfStatement = ";"
		}

		if len(input) == 0 {
			continue
		}
		if input == ":q" || input == ":quit" {
			os.Exit(0)
		}
		if input == "clear" {
			clearTerminal()
			continue
		}
		if strings.Contains(input, endOfStatement) {
			break
		} else {
			lineContinuation = true
		}

	}
	return []rune(finalInput)

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

func lastChar(input string) string {
	return string(input[len(input)-1])
}
