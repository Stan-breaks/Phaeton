package tokenize

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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
	identifier := ""
	identifierCount := 0
	for i := 0; i < len(fileContents); i++ {
		if comment > 0 && rune(fileContents[i]) != '\n' {
			continue
		}
		if !unicode.IsDigit(rune(fileContents[i])) && numberCount == 1 && rune(fileContents[i]) != utils.DOT {
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
		switch rune(fileContents[i]) {
		case utils.LEFT_PAREN:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				fmt.Println("LEFT_PAREN ( null")
			}
		case utils.RIGHT_PAREN:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					fmt.Fprintf(os.Stdout, "IDENTIFIER %s null\n", identifier)
					identifier = ""
					identifierCount = 0
				}
				fmt.Println("RIGHT_PAREN ) null")
			}
		case utils.LEFT_BRACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				fmt.Println("LEFT_BRACE { null")
			}
		case utils.RIGHT_BRACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					fmt.Fprintf(os.Stdout, "IDENTIFIER %s null\n", identifier)
					identifier = ""
					identifierCount = 0
				}
				fmt.Println("RIGHT_BRACE } null")
			}
		case utils.STAR:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {

				fmt.Println("STAR * null")
			}
		case utils.DOT:
			if numberCount == 1 {
				numberString += "."
			} else {
				if stringCount == 0 {
					fmt.Println("DOT . null")
				} else {
					stringVariable += string(rune(fileContents[i]))
				}
			}
		case utils.COMMA:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				fmt.Println("COMMA , null")
			}
		case utils.PLUS:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				fmt.Println("PLUS + null")
			}
		case utils.MINUS:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				fmt.Println("MINUS - null")
			}
		case utils.SEMICOLON:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				fmt.Println("SEMICOLON ; null")
			}
		case utils.LESS:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						fmt.Println("LESS < null")
					} else {
						switch fileContents[i+1] {
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
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						fmt.Println("GREATER > null")
					} else {
						switch fileContents[i+1] {
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
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						fmt.Println("BANG ! null")
					} else {
						switch fileContents[i+1] {
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
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						fmt.Println("EQUAL = null")
					} else {
						switch fileContents[i+1] {
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
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						fmt.Println("SLASH / null")
					} else {
						switch fileContents[i+1] {
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
			if identifierCount == 1 {
				fmt.Fprintf(os.Stdout, "IDENTIFIER %s null\n", identifier)
				identifier = ""
				identifierCount = 0
			}
		case ' ':
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					fmt.Fprintf(os.Stdout, "IDENTIFIER %s null\n", identifier)
					identifier = ""
					identifierCount = 0
				}
			}
		case '\t', '\b':
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
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
		case '@', '#', '%', '^', '&', '$':
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, rune(fileContents[i]))
				errnum = 1
				numberCount = 0
			}
		default:
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if unicode.IsDigit(rune(fileContents[i])) && identifierCount == 0 {
					if numberCount == 0 {
						numberString = ""
						numberCount = 1
						numberString += strconv.Itoa(int(rune(fileContents[i]) - '0'))
					} else {
						numberString += strconv.Itoa(int(rune(fileContents[i]) - '0'))
					}
				} else {
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
					} else {
						if identifierCount == 0 {
							switch fileContents[i] {
							case 'a':
								if fileLenght-2 > i {
									if strings.Contains(fileContents[i:i+3], "and") {
										fmt.Println("AND and null")
										if fileLenght-2 == i {
											break
										}
										i += 2
										continue
									}
								}
							case 'c':
								if fileLenght-4 > i {
									if strings.Contains(fileContents[i:i+5], "class") {
										fmt.Println("CLASS class null")
										if fileLenght-4 == i {
											break
										}
										i += 4
										continue

									}
								}
							case 'e':
								if fileLenght-3 > i {
									if strings.Contains(fileContents[i:i+4], "else") {
										fmt.Println("ELSE else null")
										if fileLenght-3 == i {
											break
										}
										i += 3
										continue
									}
								}
							case 'f':
								if fileLenght-2 > i {
									if strings.Contains(fileContents[i:i+3], "for") {
										fmt.Println("FOR for null")
										if fileLenght-2 == i {
											break
										}
										i += 2
										continue
									}
									if strings.Contains(fileContents[i:i+3], "fun") {
										fmt.Println("FUN fun null")
										if fileLenght-2 == i {
											break
										}
										i += 2
										continue
									}
								}

								if fileLenght-4 > i {
									if strings.Contains(fileContents[i:i+5], "false") {
										fmt.Println("FALSE false null")
										if fileLenght-4 == i {
											break
										}
										i += 4
										continue
									}
								}
							case 'i':
								if fileLenght-1 > i {
									if strings.Contains(fileContents[i:i+2], "if") {
										fmt.Println("IF if null")
										if fileLenght-1 == i {
											break
										}
										i += 1
										continue
									}
								}
							case 'n':
								if fileLenght-2 > i {
									if strings.Contains(fileContents[i:i+3], "nil") {
										fmt.Println("NIL nil null")
										if fileLenght-2 == i {
											break
										}
										i += 2
										continue
									}
								}
							case 'o':
								if fileLenght-1 > i {
									if strings.Contains(fileContents[i:i+2], "or") {
										fmt.Println("OR or null")
										if fileLenght-1 == i {
											break
										}
										i += 1
										continue
									}
								}
							case 'p':
								if fileLenght-4 > i {
									if strings.Contains(fileContents[i:i+5], "print") {
										fmt.Println("PRINT print null")
										if fileLenght-4 == i {
											break
										}
										i += 4
										continue
									}
								}
							case 'r':
								if fileLenght-5 > i {
									if strings.Contains(fileContents[i:i+6], "return") {
										fmt.Println("RETURN return null")
										if fileLenght-5 == i {
											break
										}
										i += 5
										continue
									}
								}
							case 's':
								if fileLenght-4 > i {
									if strings.Contains(fileContents[i:i+5], "super") {
										fmt.Println("SUPER super null")
										if fileLenght-4 == i {
											break
										}
										i += 4
										continue
									}
								}
							case 't':
								if fileLenght-3 > i {
									if strings.Contains(fileContents[i:i+4], "true") {
										fmt.Println("TRUE true null")
										if fileLenght-3 == i {
											break
										}
										i += 3
										continue
									}
									if strings.Contains(fileContents[i:i+4], "this") {
										fmt.Println("THIS this null")
										if fileLenght-3 == i {
											break
										}
										i += 3
										continue
									}
								}
							case 'v':
								if fileLenght-2 > i {
									if strings.Contains(fileContents[i:i+3], "var") {
										fmt.Println("VAR var null")
										if fileLenght-2 == i {
											break
										}
										i += 2
										continue
									}
								}
							case 'w':
								if fileLenght-4 > i {
									if strings.Contains(fileContents[i:i+5], "while") {
										fmt.Println("WHILE while null")
										if fileLenght-4 == i {
											break
										}
										i += 4
										continue
									}
								}

							}
						}
						identifier += string(rune(fileContents[i]))
						identifierCount = 1

					}
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
		if identifierCount == 1 {
			fmt.Fprintf(os.Stdout, "IDENTIFIER %s null\n", identifier)
		}
		fmt.Println("EOF  null")
	}

}
