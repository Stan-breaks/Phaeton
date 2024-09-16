package tokenize

import (
	"fmt"
	"os"

	"github.com/Stan-breaks/app/utils"
)

func Tokenize(fileContents string, fileLenght int) {
	line := 1
	errnum := 0
	skipCount := 0
	stringCount := 0
	stringVariable := ""
	comment := 0
	for index, item := range fileContents {
		if comment > 0 && item != '\n' {
			continue
		}
		switch item {
		case utils.LEFT_PAREN:
			fmt.Println("LEFT_PAREN ( null")
		case utils.RIGHT_PAREN:
			fmt.Println("RIGHT_PAREN ) null")
		case utils.LEFT_BRACE:
			fmt.Println("LEFT_BRACE { null")
		case utils.RIGHT_BRACE:
			fmt.Println("RIGHT_BRACE } null")
		case utils.STAR:
			fmt.Println("STAR * null")
		case utils.DOT:
			fmt.Println("DOT . null")
		case utils.COMMA:
			fmt.Println("COMMA , null")
		case utils.PLUS:
			fmt.Println("PLUS + null")
		case utils.MINUS:
			fmt.Println("MINUS - null")
		case utils.SEMICOLON:
			fmt.Println("SEMICOLON ; null")
		case utils.LESS:
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
						case byte(utils.EQUAL):
							fmt.Println("LESS_EQUAL <= null")
							skipCount = 1
						default:
							fmt.Println("LESS < null")
						}
					}
				}
			}

		case utils.GREATER:
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
						case byte(utils.EQUAL):
							fmt.Println("GREATER_EQUAL >= null")
							skipCount = 1
						default:
							fmt.Println("GREATER > null")
						}
					}
				}
			}

		case utils.BANG:
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
						case byte(utils.EQUAL):
							fmt.Println("BANG_EQUAL != null")
							skipCount = 1
						default:
							fmt.Println("BANG ! null")
						}
					}
				}
			}

		case utils.EQUAL:
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
						case byte(utils.EQUAL):
							fmt.Println("EQUAL_EQUAL == null")
							skipCount = 1
						default:
							fmt.Println("EQUAL = null")
						}
					}
				}
			}

		case utils.SLASH:
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
						case byte(utils.SLASH):
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
		case utils.QUOTE:
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
