package utils

import (
	"strings"
)

func IsParethesizedExpr(tokens []string) bool {
	if len(tokens) < 2 {
		return false
	}
	firstToken := strings.Split(tokens[0], " ")
	lastToken := strings.Split(tokens[len(tokens)-1], " ")
	return firstToken[0] == "LEFT_PAREN" && lastToken[0] == "RIGHT_PAREN"
}

func IsUnaryExpr(tokens []string) bool {
	if len(tokens) < 2 {
		return false
	}
	firstToken := strings.Split(tokens[0], " ")
	return (firstToken[0] == "BANG" || firstToken[0] == "MINUS" || firstToken[0] == "PLUS")
}

func IsBinaryExpression(tokens []string) bool {
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

func isInvalidToken(token string) bool {
	invalidPrefixes := []string{
		"STRING",
		"BANG",
		"TRUE",
		"FALSE",
		"IDENTIFIER",
		"LEFT_BRACE",
		"RIGHT_BRACE",
	}

	for _, prefix := range invalidPrefixes {
		if strings.HasPrefix(token, prefix) {
			return true
		}
	}
	return false
}
func IsSingleBinary(tokens []string) bool {
	numCount := 0
	for _, token := range tokens {
		if strings.HasPrefix(token, "NUMBER") {
			numCount++
		}
	}
	return numCount < 3
}

func isValidOperand(tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	if len(tokens) == 1 {
		return strings.HasPrefix(tokens[0], "NUMBER")
	}
	return IsUnaryExpr(tokens)
}

func Isoperator(token string) bool {
	operators := []string{"PLUS", "MINUS", "STAR", "SLASH"}
	for _, op := range operators {
		if strings.HasPrefix(token, op) {
			return true
		}
	}
	return false
}
