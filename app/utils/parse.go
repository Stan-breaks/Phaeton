package utils

import "strings"

func IsParethesizedExpr(tokens []string) bool {
	if len(tokens) < 2 {
		return false
	}
	firstToken := strings.Split(tokens[0], " ")
	lastToken := strings.Split(tokens[len(tokens)-1], " ")
	return firstToken[0] == "LEFT_PAREN" && lastToken[0] == "RIGHT_PAREN"
}

func IsUraryExpr(tokens []string) bool {
	firstToken := strings.Split(tokens[0], " ")
	return (firstToken[0] == "BANG" || firstToken[0] == "MINUS") && len(tokens) > 1
}

func IsBinaryExpression(tokens []string) bool {
	if len(tokens)%2 == 0 {
		return false
	}
	if !strings.HasPrefix(tokens[0], "NUMBER") {
		return false
	}

	for i := 0; i < len(tokens); i++ {
		if i%2 == 0 {
			if !strings.HasPrefix(tokens[i], "NUMBER") {
				return false
			}
		}
	}
	return true
}
