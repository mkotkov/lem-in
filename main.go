package main

import (
	"fmt"
	"os"

	"lem-in/internal"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Program takes argument!")
		os.Exit(1)
	}
	argument := os.Args[1]

	internal.RunProgram(argument, true)
}
