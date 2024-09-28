package main

import (
	"fmt"
	"github.com/Stan-breaks/app/tokenize"
	"os"
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
	tokens := tokenize.Tokenize(fileContents, fileLenght)
	switch command {
	case "tokenize":
		for _, token := range tokens.Success {
			fmt.Println(token)
		}
		fmt.Println("EOF null")
		if len(tokens.Errors) != 0 {
			for _, err := range tokens.Errors {
				fmt.Println(err)
			}
			os.Exit(65)
		}
	case "evalute":

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
