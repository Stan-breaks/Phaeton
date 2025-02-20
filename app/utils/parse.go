package utils

import (
	"github.com/Stan-breaks/app/models"
)

func IsParethesizedExpr(tokens []models.Token) bool {
	lenght := len(tokens)
	if lenght < 2 {
		return false
	}
	firstToken := tokens[0]
	lastToken := tokens[len(tokens)-1]
	if firstToken.Type != models.LEFT_PAREN || lastToken.Type != models.RIGHT_PAREN {
		return false
	}
	level := 1
	for i := 1; i < lenght-1; i++ {
		switch tokens[i].Type {
		case models.LEFT_PAREN:
			level++
		case models.RIGHT_PAREN:
			level--
		}
		if level == 0 {
			return false
		}
	}
	level--
	return level == 0
}

func IsUnaryExpr(tokens []models.Token) bool {
	firstToken := tokens[0]
	if firstToken.Type != models.BANG && firstToken.Type != models.MINUS && firstToken.Type != models.PLUS {
		return false
	}
	operandCount := 0
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Type {
		case models.NUMBER, models.STRING, models.TRUE, models.FALSE, models.IDENTIFIER:
			operandCount++
			continue
		case models.LEFT_PAREN:
			operandCount++
			parenCount := 1
			for j := i + 1; j < len(tokens); j++ {
				if token.Type == models.LEFT_PAREN {
					parenCount++
				} else if token.Type == models.RIGHT_PAREN {
					parenCount--
					if parenCount == 0 {
						i = j
						break
					}
				}
			}
		}
	}
	return operandCount == 1
}

func IsBinaryExpression(tokens []models.Token) bool {
	lenght := len(tokens)
	if lenght < 3 {
		return false
	}
	operator := -1
	for i := 1; i < lenght; i++ {
		if Isoperator(tokens[i]) {
			operator = i
		}
		if isInvalidToken(tokens[i]) {
			return false
		}
	}
	return operator != -1
}

func isInvalidToken(token models.Token) bool {
	invalidPrefixes := []models.TokenType{
		models.LEFT_BRACE,
		models.RIGHT_BRACE,
	}
	for _, prefix := range invalidPrefixes {
		if token.Type == prefix {
			return true
		}
	}
	return false
}

func IsSingleBinary(tokens []models.Token) bool {
	operandCount := 0
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Type {
		case models.NUMBER, models.STRING, models.IDENTIFIER, models.FALSE, models.TRUE, models.FUN, models.NIL:
			operandCount++
			continue
		case models.LEFT_PAREN:
			operandCount++
			parenCount := 1
			for j := i + 1; j < len(tokens); j++ {
				if token.Type == models.LEFT_PAREN {
					parenCount++
				} else if token.Type == models.RIGHT_PAREN {
					parenCount--
					if parenCount == 0 {
						i = j
						break
					}
				}
			}
		}
	}
	return operandCount == 2
}

func Isoperator(token models.Token) bool {
	operators := []models.TokenType{models.OR, models.AND, models.PLUS, models.MINUS, models.STAR, models.SLASH, models.EQUAL_EQUAL, models.BANG_EQUAL, models.LESS, models.GREATER, models.LESS_EQUAL, models.GREATER_EQUAL}
	for _, op := range operators {
		if token.Type == op {
			return true
		}
	}
	return false
}
