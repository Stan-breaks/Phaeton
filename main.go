package main

import (
	"fmt"
	"os"
)

const (
	LEFT_PAREN  rune = '('
	RIGHT_PAREN rune = ')'
	LEFT_BRACE  rune = '{'
	RIGHT_BRACE rune = '}'
	STAR        rune = '*'
	DOT         rune = '.'
	COMMA       rune = ','
	PLUS        rune = '+'
	MINUS       rune = '-'
	SEMICOLON   rune = ';'
	EQUAL       rune = '='
	BANG        rune = '!'
	LESS        rune = '<'
	GREATER     rune = '>'
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	rawFileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	line := 1
	errnum := 0
	fileContents := string(rawFileContents)
	fileLenght := len(fileContents) - 1
	skipCount := 0
	for index, item := range fileContents {
		switch item {
		case LEFT_PAREN:
			fmt.Println("LEFT_PAREN ( null")
		case RIGHT_PAREN:
			fmt.Println("RIGHT_PAREN ) null")
		case LEFT_BRACE:
			fmt.Println("LEFT_BRACE { null")
		case RIGHT_BRACE:
			fmt.Println("RIGHT_BRACE } null")
		case STAR:
			fmt.Println("STAR * null")
		case DOT:
			fmt.Println("DOT . null")
		case COMMA:
			fmt.Println("COMMA , null")
		case PLUS:
			fmt.Println("PLUS + null")
		case MINUS:
			fmt.Println("MINUS - null")
		case SEMICOLON:
			fmt.Println("SEMICOLON ; null")
		case LESS:
			if skipCount == 1 {
				skipCount = 0
				continue
			} else {
				if index == fileLenght {
					fmt.Println("LESS < null")
				} else {
					switch fileContents[index+1] {
					case byte(EQUAL):
						fmt.Println("LESS_EQUAL <= null")
						skipCount = 1
					default:
						fmt.Println("LESS < null")
					}
				}
			}
		case GREATER:
			if skipCount == 1 {
				skipCount = 0
				continue
			} else {
				if index == fileLenght {
					fmt.Println("GREATER > null")
				} else {
					switch fileContents[index+1] {
					case byte(EQUAL):
						fmt.Println("GREATER_EQUAL >= null")
						skipCount = 1
					default:
						fmt.Println("GREATER > null")
					}
				}
			}
		case BANG:
			if skipCount == 1 {
				skipCount = 0
				continue
			} else {
				if index == fileLenght {
					fmt.Println("BANG ! null")
				} else {
					switch fileContents[index+1] {
					case byte(EQUAL):
						fmt.Println("BANG_EQUAL != null")
						skipCount = 1
					default:
						fmt.Println("BANG ! null")
					}
				}
			}
		case EQUAL:
			if skipCount == 1 {
				skipCount = 0
				continue
			} else {
				if index == fileLenght {
					fmt.Println("EQUAL = null")
				} else {
					switch fileContents[index+1] {
					case byte(EQUAL):
						fmt.Println("EQUAL_EQUAL == null")
						skipCount = 1
					default:
						fmt.Println("EQUAL = null")
					}
				}
			}

		case '\n':
			line += 1
		default:
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, item)
			errnum = 1
		}
	}
	if errnum == 1 {
		fmt.Println("EOF  null")
		os.Exit(65)
	} else {
		fmt.Println("EOF  null")
	}
}
