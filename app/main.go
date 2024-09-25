package main

import (
	"fmt"
	"os"

	"github.com/Stan-breaks/app/tokenize"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	filename := os.Args[2]
	rawFileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	fileContents := string(rawFileContents)
	fileLenght := len(fileContents) - 1

	switch os.Args[1] {
	case "tokenize":
		tokenize.Tokenize(fileContents, fileLenght)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
