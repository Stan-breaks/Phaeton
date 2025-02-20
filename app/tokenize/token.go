package tokenize

import (
	"fmt"
	"math"
	"strconv"
	"unicode"

	"github.com/Stan-breaks/app/models"
	"github.com/Stan-breaks/app/utils/runes"
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
		if !unicode.IsDigit(rune(fileContents[i])) && numberCount == 1 && rune(fileContents[i]) != runes.DOT {
			if numberCount == 1 {
				number, err := strconv.Atoi(numberString)
				if err != nil {
					float, err := strconv.ParseFloat(numberString, 64)
					if err != nil {
						fmt.Println("Error parsing float:", err)
					}
					if math.Mod(float, 1.0) == 0 {
						token := models.Token{
							Type:    models.NUMBER,
							Lexem:   numberString,
							Literal: float,
							Line:    line,
						}
						tokens.Success = append(tokens.Success, token)
						numberCount = 0
					} else {
						token := models.Token{
							Type:    models.NUMBER,
							Lexem:   numberString,
							Literal: float,
							Line:    line,
						}
						tokens.Success = append(tokens.Success, token)
						numberCount = 0
					}
				} else {
					token := models.Token{
						Type:    models.NUMBER,
						Lexem:   numberString,
						Literal: number,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					numberCount = 0
				}
			}
		}
		switch rune(fileContents[i]) {
		case runes.LEFT_PAREN:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.Token{
					Type:    models.LEFT_PAREN,
					Lexem:   "(",
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case runes.RIGHT_PAREN:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.Token{
					Type:    models.RIGHT_PAREN,
					Lexem:   ")",
					Literal: nil,
					Line:    line,
				}

				tokens.Success = append(tokens.Success, token)
			}
		case runes.LEFT_BRACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.Token{
					Type:    models.LEFT_BRACE,
					Lexem:   "{",
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case runes.RIGHT_BRACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.Token{
					Type:    models.RIGHT_BRACE,
					Lexem:   "}",
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case runes.STAR:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				token := models.Token{
					Type:    models.STAR,
					Lexem:   "*",
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case runes.DOT:
			if numberCount == 1 {
				numberString += "."
			} else {
				if stringCount == 0 {
					token := models.Token{
						Type:    models.DOT,
						Lexem:   ".",
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
				} else {
					stringVariable += string(rune(fileContents[i]))
				}
			}
		case runes.COMMA:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.Token{
					Type:    models.COMMA,
					Lexem:   ",",
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case runes.PLUS:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.Token{
					Type:    models.PLUS,
					Lexem:   "+",
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case runes.MINUS:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.Token{
					Type:    models.MINUS,
					Lexem:   "-",
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case runes.SEMICOLON:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
				token := models.Token{
					Type:    models.SEMICOLON,
					Lexem:   ";",
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
			}
		case runes.LESS:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						token := models.Token{
							Type:    models.LESS,
							Lexem:   "<",
							Literal: nil,
							Line:    line,
						}
						tokens.Success = append(tokens.Success, token)
					} else {
						switch fileContents[i+1] {
						case byte(runes.EQUAL):
							token := models.Token{
								Type:    models.LESS_EQUAL,
								Lexem:   "<=",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
							skipCount = 1
						default:
							token := models.Token{
								Type:    models.LESS,
								Lexem:   "<",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
						}
					}
				}
			}

		case runes.GREATER:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						token := models.Token{
							Type:    models.GREATER,
							Lexem:   ">",
							Literal: nil,
							Line:    line,
						}
						tokens.Success = append(tokens.Success, token)
					} else {
						switch fileContents[i+1] {
						case byte(runes.EQUAL):
							token := models.Token{
								Type:    models.GREATER_EQUAL,
								Lexem:   ">=",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
							skipCount = 1
						default:
							token := models.Token{
								Type:    models.GREATER,
								Lexem:   ">",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
						}
					}
				}
			}

		case runes.BANG:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						token := models.Token{
							Type:    models.BANG,
							Lexem:   "!",
							Literal: nil,
							Line:    line,
						}
						tokens.Success = append(tokens.Success, token)
					} else {
						switch fileContents[i+1] {
						case byte(runes.EQUAL):
							token := models.Token{
								Type:    models.BANG_EQUAL,
								Lexem:   "!=",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
							skipCount = 1
						default:
							token := models.Token{
								Type:    models.BANG,
								Lexem:   "!",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
						}
					}
				}
			}

		case runes.EQUAL:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						token := models.Token{
							Type:    models.EQUAL,
							Lexem:   "=",
							Literal: nil,
							Line:    line,
						}
						tokens.Success = append(tokens.Success, token)
					} else {
						switch fileContents[i+1] {
						case byte(runes.EQUAL):
							token := models.Token{
								Type:    models.EQUAL_EQUAL,
								Lexem:   "==",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
							skipCount = 1
						default:
							token := models.Token{
								Type:    models.EQUAL,
								Lexem:   "=",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
						}
					}
				}
			}

		case runes.SLASH:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if skipCount == 1 {
					skipCount = 0
				} else {
					if i == fileLenght {
						token := models.Token{
							Type:    models.SLASH,
							Lexem:   "/",
							Literal: nil,
							Line:    line,
						}
						tokens.Success = append(tokens.Success, token)
					} else {
						switch fileContents[i+1] {
						case byte(runes.SLASH):
							if stringCount == 0 {
								comment = 1
							}
						default:
							token := models.Token{
								Type:    models.SLASH,
								Lexem:   "/",
								Literal: nil,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
						}
					}
				}
			}

		case runes.NEWLINE:
			numberCount = 0
			line += 1
			comment = 0
			if identifierCount == 1 {
				token := models.Token{
					Type:    models.IDENTIFIER,
					Lexem:   identifier,
					Literal: nil,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
				identifier = ""
				identifierCount = 0
			}
		case runes.WHITESPACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			} else {
				if identifierCount == 1 {
					token := models.Token{
						Type:    models.IDENTIFIER,
						Lexem:   identifier,
						Literal: nil,
						Line:    line,
					}
					tokens.Success = append(tokens.Success, token)
					identifier = ""
					identifierCount = 0
				}
			}
		case runes.TAB, runes.BACKSPACE:
			numberCount = 0
			if stringCount == 1 {
				stringVariable += string(rune(fileContents[i]))
			}
		case runes.QUOTE:
			numberCount = 0
			if stringCount == 1 {
				stringCount = 0
				token := models.Token{
					Type:    models.STRING,
					Lexem:   "\"" + stringVariable + "\"",
					Literal: stringVariable,
					Line:    line,
				}
				tokens.Success = append(tokens.Success, token)
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
							token := models.Token{
								Type:    models.NUMBER,
								Lexem:   numberString,
								Literal: float,
								Line:    line,
							}
							tokens.Success = append(tokens.Success, token)
							numberCount = 0

						} else {
							token := models.Token{
								Type:    models.NUMBER,
								Lexem:   numberString,
								Literal: number,
								Line:    line,
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
										token := models.TokenInfo{
											Token: "AND and null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "CLASS class null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
										if fileLenght-3 == i {
											break
										}
										i += 4
										continue

									}
								}
								if fileLenght-5 > i {
									if i+6 <= fileLenght && fileContents[i:i+7] == "clock()" {
										token := models.TokenInfo{
											Token: "FUNCTION clock null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
										if fileLenght-5 == i {
											break
										}
										i += 6
										continue
									}
								}
							case 'e':
								if fileLenght-2 > i {
									if i+3 <= fileLenght && fileContents[i:i+4] == "else" {
										token := models.TokenInfo{
											Token: "ELSE else null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "FOR for null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
										if fileLenght-1 == i {
											break
										}
										i += 2
										continue
									}
									if i+2 <= fileLenght && fileContents[i:i+3] == "fun" {
										token := models.TokenInfo{
											Token: "FUN fun null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
										if fileLenght-1 == i {
											break
										}
										i += 2
										continue
									}
								}

								if fileLenght-3 > i {
									if i+4 <= fileLenght && fileContents[i:i+5] == "false" {
										token := models.TokenInfo{
											Token: "FALSE false null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "IF if null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "NIL nil null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "OR or null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "PRINT print null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "RETURN return null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "SUPER super null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "TRUE true null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
										if fileLenght-2 == i {
											break
										}
										i += 3
										continue
									}
									if i+3 <= fileLenght && fileContents[i:i+4] == "this" {
										token := models.TokenInfo{
											Token: "THIS this null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "VAR var null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
										token := models.TokenInfo{
											Token: "WHILE while null",
											Line:  line,
										}
										tokens.Success = append(tokens.Success, token)
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
