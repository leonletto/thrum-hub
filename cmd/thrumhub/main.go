package main

import (
	"fmt"
	"os"
)

var version = "0.1.0-dev"

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "version" {
		fmt.Println("thrumhub", version)
		return
	}
	fmt.Fprintln(os.Stderr, "thrumhub: no command given (try `thrumhub version`)")
	os.Exit(2)
}
