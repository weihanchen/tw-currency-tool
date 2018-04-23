package main

import (
	"os"

	"gitlab.com/will.chen/tw-currency-tool/commands"
)

func main() {
	if err := commands.Execute(os.Args[1:]); err != nil {
		// fmt.Println(err)
		os.Exit(-1)
	}
}
