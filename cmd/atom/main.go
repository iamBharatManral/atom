package main

import (
	"os"

	filerunner "github.com/iamBharatManral/atom.git/cmd/internal/fileRunner"
	"github.com/iamBharatManral/atom.git/cmd/internal/repl"
	"github.com/iamBharatManral/atom.git/cmd/internal/util"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		repl.Start()
	} else if len(args) == 2 {
		if args[1] == "-h" {
			util.Usage()
			os.Exit(0)
		} else {
			filename := args[1]
			filerunner.Execute(filename)
		}
	}
}
