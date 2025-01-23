package utils

import (
	"strings"

	"github.com/Stan-breaks/app/models"
)

func IsParethesizedExpr(tokens []models.TokenInfo) bool {
	lenght := len(tokens)
	if lenght < 2 {
		return false
	}
	firstToken := strings.Split(tokens[0].Token, " ")
	lastToken := strings.Split(tokens[len(tokens)-1].Token, " ")
	if firstToken[0] != "LEFT_PAREN" || lastToken[0] != "RIGHT_PAREN" {
		return false
	}
	level := 1
	for i := 1; i < lenght-1; i++ {
		tokenType := strings.Split(tokens[i].Token, " ")[0]
		switch tokenType {
		case "LEFT_PAREN":
			level++
		case "RIGHT_PAREN":
			level--
		}
		if level == 0 {
			return false
		}
	}
	level--
	return level == 0
}

func IsUnaryExpr(tokens []models.TokenInfo) bool {
	firstToken := strings.Split(tokens[0].Token, " ")
	if firstToken[0] != "BANG" && firstToken[0] != "MINUS" && firstToken[0] != "PLUS" {
		return false
	}
	operandCount := 0
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if strings.HasPrefix(token.Token, "NUMBER") ||
			strings.HasPrefix(token.Token, "STRING") ||
			strings.HasPrefix(token.Token, "TRUE") ||
			strings.HasPrefix(token.Token, "FALSE") ||
			strings.HasPrefix(token.Token, "IDENTIFIER") {
			operandCount++
			continue
		}
		if strings.HasPrefix(token.Token, "LEFT_PAREN") {
			operandCount++
			parenCount := 1
			for j := i + 1; j < len(tokens); j++ {
				if strings.HasPrefix(tokens[j].Token, "LEFT_PAREN") {
					parenCount++
				} else if strings.HasPrefix(tokens[j].Token, "RIGHT_PAREN") {
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

func IsBinaryExpression(tokens []models.TokenInfo) bool {
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

func isInvalidToken(token models.TokenInfo) bool {
	invalidPrefixes := []string{
		"LEFT_BRACE",
		"RIGHT_BRACE",
	}
	for _, prefix := range invalidPrefixes {
		if strings.HasPrefix(token.Token, prefix) {
			return true
		}
	}
	return false
}

func IsSingleBinary(tokens []models.TokenInfo) bool {
	operandCount := 0
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if strings.HasPrefix(token.Token, "NUMBER") ||
			strings.HasPrefix(token.Token, "STRING") ||
			strings.HasPrefix(token.Token, "IDENTIFIER") ||
			strings.HasPrefix(token.Token, "FALSE") ||
			strings.HasPrefix(token.Token, "TRUE") ||
			strings.HasPrefix(token.Token, "NIL") {
			operandCount++
			continue
		} else if strings.HasPrefix(token.Token, "LEFT_PAREN") {
			operandCount++
			parenCount := 1
			for j := i + 1; j < len(tokens); j++ {
				if strings.HasPrefix(tokens[j].Token, "LEFT_PAREN") {
					parenCount++
				} else if strings.HasPrefix(tokens[j].Token, "RIGHT_PAREN") {
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

func Isoperator(token models.TokenInfo) bool {
	operators := []string{"OR", "PLUS", "MINUS", "STAR", "SLASH", "EQUAL_EQUAL", "BANG_EQUAL", "LESS", "GREATER", "LESS_EQUAL", "GREATER_EQUAL"}
	for _, op := range operators {
		if strings.HasPrefix(token.Token, op) {
			return true
		}
	}
	return false
}
