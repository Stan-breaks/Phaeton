package tokenize

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"unicode"

	"github.com/Stan-breaks/app/utils"
)

func Tokenize(fileContents string, fileLenght int) {
	line := 1
	errnum := 0
	skipCount := 0
	stringCount := 0
	stringVariable := ""
	numberCount := 0
	numberString := ""
	comment := 0
	for index, item := range fileContents {
		if comment > 0 && item != '\n' {
			continue
		}
		if !unicode.IsDigit(item) && numberCount == 1 && item != utils.DOT {
			if numberCount == 1 {
				number, err := strconv.Atoi(numberString)
				if err != nil {
					float, err := strconv.ParseFloat(numberString, 64)
					if err != nil {
						fmt.Println("Error parsing float:", err)
					}
					if math.Mod(float, 1.0) == 0 {
						fmt.Fprintf(os.Stdout, "NUMBER %s %.1f\n", numberString, float)
						numberCount = 0
					} else {
						fmt.Fprintf(os.Stdout, "NUMBER %s %g\n", numberString, float)
						numberCount = 0

					}
				} else {
					fmt.Fprintf(os.Stdout, "NUMBER %s %d.0\n", numberString, number)
					numberCount = 0
				}
			}
		}
		switch item {
		case utils.LEFT_PAREN:
			numberCount = 0
			fmt.Println("LEFT_PAREN ( null")
		case utils.RIGHT_PAREN:
			numberCount = 0
			fmt.Println("RIGHT_PAREN ) null")
		case utils.LEFT_BRACE:
			numberCount = 0
			fmt.Println("LEFT_BRACE { null")
		case utils.RIGHT_BRACE:
			numberCount = 0
			fmt.Println("RIGHT_BRACE } null")
		case utils.STAR:
			numberCount = 0
			fmt.Println("STAR * null")
		case utils.DOT:
			if numberCount == 1 {
				numberString += "."
			} else {
				if stringCount == 0 {
					fmt.Println("DOT . null")
				}
			}
		case utils.COMMA:
			numberCount = 0
			fmt.Println("COMMA , null")
		case utils.PLUS:
			numberCount = 0
			fmt.Println("PLUS + null")
		case utils.MINUS:
			numberCount = 0
			fmt.Println("MINUS - null")
		case utils.SEMICOLON:
			numberCount = 0
			fmt.Println("SEMICOLON ; null")
		case utils.LESS:
			numberCount = 0
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
			numberCount = 0
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
			numberCount = 0
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
			numberCount = 0
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
			numberCount = 0
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
			numberCount = 0
			line += 1
			comment = 0
		case '\t', '\b', ' ':
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(item)
			}
		case utils.QUOTE:
			numberCount = 0
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
				if unicode.IsDigit(item) {
					if numberCount == 0 {
						numberString = ""
						numberCount = 1
						numberString += strconv.Itoa(int(item - '0'))
					} else {
						numberString += strconv.Itoa(int(item - '0'))
					}
				} else {

					fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, item)
					errnum = 1
					numberCount = 0
				}
			}
		}
	}
	if stringCount == 1 {
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
	}
	if numberCount == 1 {
		number, err := strconv.Atoi(numberString)
		if err != nil {
			float, err := strconv.ParseFloat(numberString, 64)
			if err != nil {
				fmt.Println("Error parsing float:", err)
			}
			if math.Mod(float, 1.0) == 0 {
				fmt.Fprintf(os.Stdout, "NUMBER %s %.1f\n", numberString, float)
			} else {
				fmt.Fprintf(os.Stdout, "NUMBER %s %g\n", numberString, float)

			}
		} else {
			fmt.Fprintf(os.Stdout, "NUMBER %s %d.0\n", numberString, number)
		}
	}
	if errnum == 1 || stringCount == 1 {
		fmt.Println("EOF  null")
		os.Exit(65)
	} else {
		fmt.Println("EOF  null")
	}

}
