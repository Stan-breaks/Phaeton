package main

import (
	"fmt"
	"os"

	"github.com/Stan-breaks/app/interpreter"
	"github.com/Stan-breaks/app/parse"
	"github.com/Stan-breaks/app/tokenize"
)

func main() {
	if len(os.Args) < 3 {
		if os.Args[1] == "--version" {
			fmt.Println("Phaeton v0.0.0")
			os.Exit(0)
		}
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
			fmt.Println(token.Token)
		}
		fmt.Println("EOF  null")
	case "parse":
		var err []string
		value, err := parse.Parse(tokens.Success)
		tokens.Errors = append(tokens.Errors, err...)
		result := value.String()
		fmt.Printf("%s\n", result)
	case "run":
		err := interpreter.Interprete(tokens.Success)
		if err != nil {
			tokens.Errors = append(tokens.Errors, err.Error())
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
	if len(tokens.Errors) != 0 {
		for _, err := range tokens.Errors {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(65)
	}
}
