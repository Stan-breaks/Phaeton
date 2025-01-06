package tokenize

import (
	"fmt"
	"go/token"
	"math"
	"strconv"
	"unicode"

	"github.com/Stan-breaks/app/models"
	"github.com/Stan-breaks/app/utils"
	"golang.org/x/mod/module"
)

func Tokenize(fileContents string, fileLenght int) models.Tokens {
	tokens := models.Tokens{}
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
						token := models.TokenInfo{
							Token: fmt.Sprintf("NUMBER %s %.1f", numberString, float),
							Line:  line,
						}
						tokens.Success = append(tokens.Success, token)
						numberCount = 0
					} else {
						token := models.TokenInfo{
							Token: fmt.Sprintf("NUMBER %s %g", numberString, float),
							Line:  line,
						}
						tokens.Success = append(tokens.Success, token)
						numberCount = 0
					}
				} else {
					token := models.TokenInfo{
						Token: fmt.Sprintf("NUMBER %s %d.0", numberString, number),
						Line:  line,
					}
					tokens.Success = append(tokens.Success, token)
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
				token := models.TokenInfo{
					Token: "LEFT_PAREN ( null",
					Line:  line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case utils.RIGHT_PAREN:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.TokenInfo{
						Token: fmt.Sprintf("IDENTIFIER %s null", identifier),
						Line:  line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.TokenInfo{
					Token: "RIGHT_PAREN ) null",
					Line:  line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case utils.LEFT_BRACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				tokens.Success = append(tokens.Success, "LEFT_BRACE { null")
			}
		case utils.RIGHT_BRACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.TokenInfo{
						Token: fmt.Sprintf("IDENTIFIER %s null", identifier),
						Line:  line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				tokens.Success = append(tokens.Success, "RIGHT_BRACE } null")
			}
		case utils.STAR:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				tokens.Success = append(tokens.Success, "STAR * null")
			}
		case utils.DOT:
			if numberCount == 1 {
				numberString += "."
			} else {
				if stringCount == 0 {
					tokens.Success = append(tokens.Success, "DOT . null")
				} else {
					stringVariable += string(rune(fileContents[i]))
				}
			}
		case utils.COMMA:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				tokens.Success = append(tokens.Success, "COMMA , null")
			}
		case utils.PLUS:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				tokens.Success = append(tokens.Success, "PLUS + null")
			}
		case utils.MINUS:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				tokens.Success = append(tokens.Success, "MINUS - null")
			}
		case utils.SEMICOLON:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				tokens.Success = append(tokens.Success, "SEMICOLON ; null")
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
						tokens.Success = append(tokens.Success, "LESS < null")
					} else {
						switch fileContents[i+1] {
						case byte(utils.EQUAL):
							tokens.Success = append(tokens.Success, "LESS_EQUAL <= null")
							skipCount = 1
						default:
							tokens.Success = append(tokens.Success, "LESS < null")
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
						tokens.Success = append(tokens.Success, "GREATER > null")
					} else {
						switch fileContents[i+1] {
						case byte(utils.EQUAL):
							tokens.Success = append(tokens.Success, "GREATER_EQUAL >= null")
							skipCount = 1
						default:
							tokens.Success = append(tokens.Success, "GREATER > null")
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
						tokens.Success = append(tokens.Success, "BANG ! null")
					} else {
						switch fileContents[i+1] {
						case byte(utils.EQUAL):
							tokens.Success = append(tokens.Success, "BANG_EQUAL != null")
							skipCount = 1
						default:
							tokens.Success = append(tokens.Success, "BANG ! null")
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
						tokens.Success = append(tokens.Success, "EQUAL = null")
					} else {
						switch fileContents[i+1] {
						case byte(utils.EQUAL):
							tokens.Success = append(tokens.Success, "EQUAL_EQUAL == null")
							skipCount = 1
						default:
							tokens.Success = append(tokens.Success, "EQUAL = null")
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
						tokens.Success = append(tokens.Success, "SLASH / null")
					} else {
						switch fileContents[i+1] {
						case byte(utils.SLASH):
							if stringCount == 0 {
								comment = 1
							}
						default:
							tokens.Success = append(tokens.Success, "SLASH / null")
						}
					}
				}
			}

		case utils.NEWLINE:
			numberCount = 0
			line += 1
			comment = 0
			if identifierCount == 1 {
				token := models.TokenInfo{
					Token: fmt.Sprintf("IDENTIFIER %s null", identifier),
					Line:  line,
				}
				tokens.Success = append(tokens.Success, token)
				identifier = ""
				identifierCount = 0
			}
		case utils.WHITESPACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.TokenInfo{
						Token: fmt.Sprintf("IDENTIFIER %s null", identifier),
						Line:  line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
			}
		case utils.TAB, utils.BACKSPACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			}
		case utils.QUOTE:
			numberCount = 0
			if stringCount == 1 {
				stringCount = 0
				tokens.Success = append(tokens.Success, "STRING \""+stringVariable+"\" "+stringVariable)
			} else {
				stringCount = 1
				stringVariable = ""
			}
		case '@', '#', '%', '^', '&', '$':
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				tokens.Errors = append(tokens.Errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c\n", line, rune(fileContents[i])))
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
								token := models.TokenInfo{
									Token: fmt.Sprintf("NUMBER %s %.1f", numberString, float),
									Line:  line,
								}
								tokens.Success = append(tokens.Success, token)
								numberCount = 0
							} else {
								token := models.TokenInfo{
									Token: fmt.Sprintf("NUMBER %s %g", numberString, float),
									Line:  line,
								}
								tokens.Success = append(tokens.Success, token)
								numberCount = 0
							}
						} else {
							token := models.TokenInfo{
								Token: fmt.Sprintf("NUMBER %s %d.0", numberString, number),
								Line:  line,
							}
							tokens.Success = append(tokens.Success, token)
							numberCount = 0
						}
					} else {
						if identifierCount == 0 {
							switch fileContents[i] {
							case 'a':
								if fileLenght-1 > i {
									if i+2 <= fileLenght && fileContents[i:i+3] == "and" {
										tokens.Success = append(tokens.Success, "AND and null")
										if fileLenght-1 == i {
											break
										}
										i += 2
										continue
									}
								}
							case 'c':
								if fileLenght-3 > i {
									if i+4 <= fileLenght && fileContents[i:i+5] == "class" {
										tokens.Success = append(tokens.Success, "CLASS class null")
										if fileLenght-3 == i {
											break
										}
										i += 4
										continue

									}
								}
							case 'e':
								if fileLenght-2 > i {
									if i+3 <= fileLenght && fileContents[i:i+4] == "else" {
										tokens.Success = append(tokens.Success, "ELSE else null")
										if fileLenght-2 == i {
											break
										}
										i += 3
										continue
									}
								}
							case 'f':
								if fileLenght-1 > i {
									if i+2 <= fileLenght && fileContents[i:i+3] == "for" {
										tokens.Success = append(tokens.Success, "FOR for null")
										if fileLenght-1 == i {
											break
										}
										i += 2
										continue
									}
									if i+2 <= fileLenght && fileContents[i:i+3] == "fun" {
										tokens.Success = append(tokens.Success, "FUN fun null")
										if fileLenght-1 == i {
											break
										}
										i += 2
										continue
									}
								}

								if fileLenght-3 > i {
									if i+4 <= fileLenght && fileContents[i:i+5] == "false" {
										tokens.Success = append(tokens.Success, "FALSE false null")
										if fileLenght-3 == i {
											break
										}
										i += 4
										continue
									}
								}
							case 'i':
								if fileLenght > i {
									if i+1 <= fileLenght && fileContents[i:i+2] == "if" {
										tokens.Success = append(tokens.Success, "IF if null")
										if fileLenght == i {
											break
										}
										i += 1
										continue
									}
								}
							case 'n':
								if fileLenght-1 > i {
									if i+2 <= fileLenght && fileContents[i:i+3] == "nil" {
										tokens.Success = append(tokens.Success, "NIL nil null")
										if fileLenght-1 == i {
											break
										}
										i += 2
										continue
									}
								}
							case 'o':
								if fileLenght > i {
									if i+1 <= fileLenght && fileContents[i:i+2] == "or" {
										tokens.Success = append(tokens.Success, "OR or null")
										if fileLenght == i {
											break
										}
										i += 1
										continue
									}
								}
							case 'p':
								if fileLenght-3 > i {
									if i+4 <= fileLenght && fileContents[i:i+5] == "print" {
										tokens.Success = append(tokens.Success, "PRINT print null")
										if fileLenght-3 == i {
											break
										}
										i += 4
										continue
									}
								}
							case 'r':
								if fileLenght-4 > i {
									if i+5 <= fileLenght && fileContents[i:i+6] == "return" {
										tokens.Success = append(tokens.Success, "RETURN return null")
										if fileLenght-4 == i {
											break
										}
										i += 5
										continue
									}
								}
							case 's':
								if fileLenght-3 > i {
									if i+4 <= fileLenght && fileContents[i:i+5] == "super" {
										tokens.Success = append(tokens.Success, "SUPER super null")
										if fileLenght-3 == i {
											break
										}
										i += 4
										continue
									}
								}
							case 't':
								if fileLenght-2 > i {
									if i+3 <= fileLenght && fileContents[i:i+4] == "true" {
										tokens.Success = append(tokens.Success, "TRUE true null")
										if fileLenght-2 == i {
											break
										}
										i += 3
										continue
									}
									if i+3 <= fileLenght && fileContents[i:i+4] == "this" {
										tokens.Success = append(tokens.Success, "THIS this null")
										if fileLenght-2 == i {
											break
										}
										i += 3
										continue
									}
								}
							case 'v':
								if fileLenght-1 > i {
									if i+2 <= fileLenght && fileContents[i:i+3] == "var" {
										tokens.Success = append(tokens.Success, "VAR var null")
										if fileLenght-1 == i {
											break
										}
										i += 2
										continue
									}
								}
							case 'w':
								if fileLenght-3 > i {
									if i+4 <= fileLenght && fileContents[i:i+5] == "while" {
										tokens.Success = append(tokens.Success, "WHILE while null")
										if fileLenght-3 == i {
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
		tokens.Errors = append(tokens.Errors, fmt.Sprintf("[line %d] Error: Unterminated string.", line))
	}
	if numberCount == 1 {
		number, err := strconv.Atoi(numberString)
		if err != nil {
			float, err := strconv.ParseFloat(numberString, 64)
			if err != nil {
				fmt.Println("Error parsing float:", err)
			}
			if math.Mod(float, 1.0) == 0 {
				token := models.TokenInfo{
					Token: fmt.Sprintf("NUMBER %s %.1f", numberString, float),
					Line:  line,
				}
				tokens.Success = append(tokens.Success, token)
			} else {
				token := models.TokenInfo{
					Token: fmt.Sprintf("NUMBER %s %g", numberString, float),
					Line:  line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		} else {
			token := models.TokenInfo{
				Token: fmt.Sprintf("NUMBER %s %d.0", numberString, number),
				Line:  line,
			}
			tokens.Success = append(tokens.Success, token)
		}
	}
	if errnum != 1 || stringCount != 1 {
		if identifierCount == 1 {
			token := models.TokenInfo{
				Token: fmt.Sprintf("IDENTIFIER %s null", identifier),
				Line:  line,
			}
			tokens.Success = append(tokens.Success, token)
		}
	}
	return tokens
}
