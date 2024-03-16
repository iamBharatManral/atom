package repl

import (
	"fmt"
	"os/user"

	"github.com/iamBharatManral/atom.git/cmd/internal/util"
)

func Start() {
	util.Banner()
	message()
}

func message() {
	var username string
	currentUser, err := user.Current()
	if err != nil {
		username = "STRANGER"
	} else {
		username = currentUser.Username
	}
	fmt.Printf("Welcome %s to the beginning of the language universe", username)
}
