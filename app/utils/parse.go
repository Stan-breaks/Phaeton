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
	if lenght < 3 || lenght > 5 {
		return false
	}

	operator := -1
	for i := 1; i < len(tokens); i++ {
		if isOperator(tokens[i]) {
			operator = i
			break
		}
	}
	leftOperand := tokens[:operator]
	rightOperand := tokens[operator+1:]
	return isValidOperand(leftOperand) && isValidOperand(rightOperand)
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

func isOperator(token string) bool {
	operators := []string{"PLUS", "MINUS", "STAR", "SLASH"}
	for _, op := range operators {
		if strings.HasPrefix(token, op) {
			return true
		}
	}
	return false
}
