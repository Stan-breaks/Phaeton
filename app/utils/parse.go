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
	return len(tokens) == 3 && strings.HasPrefix(tokens[0], "NUMBER")
}
