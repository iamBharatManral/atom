package util

import "fmt"

func Banner() {
	fmt.Println("\n\t  __   ____   __   _  _ ")
	fmt.Println("\t / _\\ (_  _) /  \\ ( \\/ )")
	fmt.Println("\t/    \\  )(  (  O )/ \\/ \\")
	fmt.Println("\t\\_/\\_/ (__)  \\__/ \\_)(_/")
	fmt.Println()
}

func Usage() {
	fmt.Println("\nThere are two ways to explore atom: 🌖")
	fmt.Println("\t1. atom <enter>: will open atom interpreter")
	fmt.Println("\t2. atom <file.om>: will execute the file")
	fmt.Println("Enjoy!")
	fmt.Println("")
}
