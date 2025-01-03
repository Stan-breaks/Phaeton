package utils

import (
	"strings"
)

func IsParethesizedExpr(tokens []string) bool {
	lenght := len(tokens)
	if lenght < 2 {
		return false
	}
	firstToken := strings.Split(tokens[0], " ")
	lastToken := strings.Split(tokens[len(tokens)-1], " ")
	if firstToken[0] != "LEFT_PAREN" || lastToken[0] != "RIGHT_PAREN" {
		return false
	}
	level := 1
	for i := 1; i < lenght-1; i++ {
		tokenType := strings.Split(tokens[i], " ")[0]
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

func IsUnaryExpr(tokens []string) bool {
	if len(tokens) > 2 {
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
	stringCount := 0
	for _, token := range tokens {
		if strings.HasPrefix(token, "STRING") {
			stringCount++
		}
	}
	return numCount == 2 || stringCount == 2
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
