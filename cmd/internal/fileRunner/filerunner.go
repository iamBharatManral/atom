package filerunner

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/iamBharatManral/atom.git/cmd/internal/env"
	"github.com/iamBharatManral/atom.git/cmd/internal/interpreter"
	"github.com/iamBharatManral/atom.git/cmd/internal/lexer"
	"github.com/iamBharatManral/atom.git/cmd/internal/parser"
)

func Execute(filename string) {
	fileInfo := strings.Split(filename, ".")
	if fileInfo[1] != "om" {
		log.Printf("error: wrong filetype, %s is not .om file", filename)
		os.Exit(1)
	}
	input, err := os.ReadFile(filename)
	if err != nil && !os.IsExist(err) {
		log.Printf("error: file %s, does not exists", filename)
		os.Exit(1)
	}
	lexer := lexer.New([]rune(string(input)))
	parser := parser.New(lexer)
	program := parser.Parse()
	env := env.New(nil)
	result := interpreter.Eval(program, env)
	if result.Type == "error" {
		log.Fatal(result.Value)
	} else if result.Type == "" {
		return
	}
	fmt.Println(result.Value)
}
