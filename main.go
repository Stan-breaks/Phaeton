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
	SLASH       rune = '/'
	QUOTE       rune = '"'
)

func main() {
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
	stringCount := 0
	stringVariable := ""
	comment := 0
	for index, item := range fileContents {
		if comment > 0 && item != '\n' {
			continue
		}
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
			if stringCount == 1 {
				stringVariable += string(item)
			} else {
				if skipCount == 1 {
					skipCount = 0
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
			}

		case GREATER:
			if stringCount == 1 {
				stringVariable += string(item)
			} else {
				if skipCount == 1 {
					skipCount = 0
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
			}

		case BANG:
			if stringCount == 1 {
				stringVariable += string(item)
			} else {
				if skipCount == 1 {
					skipCount = 0
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
			}

		case EQUAL:
			if stringCount == 1 {
				stringVariable += string(item)
			} else {
				if skipCount == 1 {
					skipCount = 0
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
			}

		case SLASH:
			if stringCount == 1 {
				stringVariable += string(item)
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if index == fileLenght {
						fmt.Println("SLASH / null")
					} else {
						switch fileContents[index+1] {
						case byte(SLASH):
							if stringCount == 0 {
								comment = 1
							}
						default:
							continue
						}
					}
				}
			}

		case '\n':
			line += 1
			comment = 0
		case '\t', '\b', ' ':
			if stringCount == 1 {
				stringVariable += string(item)
			}
		case QUOTE:
			if stringCount == 1 {
				stringCount = 0
				fmt.Println("STRING \"" + stringVariable + "\" " + stringVariable)
			} else {
				stringCount = 1
				stringVariable = ""
			}
		default:
			if stringCount == 1 {
				stringVariable += string(item)
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, item)
				errnum = 1
			}
		}
	}
	if stringCount == 1 {
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
	}
	if errnum == 1 || stringCount == 1 {
		fmt.Println("EOF  null")
		os.Exit(65)
	} else {
		fmt.Println("EOF  null")
	}
}
