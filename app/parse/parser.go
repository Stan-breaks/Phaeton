package parse

import (
	"strconv"
	"strings"

	"github.com/Stan-breaks/app/models"
	"github.com/Stan-breaks/app/utils"
)

func Parse(tokens models.Tokens) models.Node {
	if len(tokens.Success) == 1 {
		splitToken := strings.Split(tokens.Success[0], " ")
		return parsevalue(splitToken)
	}
	if utils.IsParethesizedExpr(tokens.Success) {
		return parseParrenthesisExpr(tokens.Success)
	}
	if utils.IsBinaryExpression(tokens.Success) {
		return parseBinaryExpr(tokens.Success)
	}
	if utils.IsUraryExpr(tokens.Success) {
		return parseUraryExpr(tokens.Success)
	}
	return models.NilNode{}
}

func parseBinaryExpr(tokens []string) models.Node {
	if len(tokens) <= 3 {
		return parseSingleBinaryExpr(tokens)
	} else {
		return parseMultipleBinaryExpr(tokens)
	}
}

func parseSingleBinaryExpr(tokens []string) models.Node {
	var left, right models.Node
	op := ""
	for index, item := range tokens {
		splitToken := strings.Split(item, " ")
		switch index {
		case 0:
			left = parsevalue(splitToken)
		case 1:
			op = parseOperator(splitToken)
		case 2:
			right = parsevalue(splitToken)
		}
	}
	return models.BinaryNode{
		Left:  left,
		Op:    op,
		Right: right,
	}
}

func parseMultipleBinaryExpr(tokens []string) models.Node {
	var left, right models.Node
	firstBinary := tokens[0:3]
	left = parseSingleBinaryExpr(firstBinary)
	splitToken := strings.Split(tokens[3], " ")
	op := parseOperator(splitToken)
	remainingTokens := tokens[4:]
	if len(remainingTokens) == 1 {
		right = parsevalue(strings.Split(remainingTokens[0], " "))
	} else {
		right = parseBinaryExpr(remainingTokens)
	}
	binary := models.BinaryNode{
		Left:  left,
		Op:    op,
		Right: right,
	}

	return models.StringNode{
		Value: binary.String(),
	}
}

func parseOperator(splitToken []string) string {
	switch splitToken[0] {
	case "PLUS", "MINUS", "STAR", "SLASH", "EQUAL_EQUAL", "LESS", "AND", "OR":
		return splitToken[1]
	default:
		return ""
	}
}

func parsevalue(splitToken []string) models.Node {
	switch splitToken[0] {
	case "NUMBER":
		num, _ := strconv.ParseFloat(splitToken[1], 32)
		floatnum := float32(num)
		return models.NumberNode{Value: floatnum}
	case "TRUE":
		return models.BooleanNode{Value: true}
	case "FALSE":
		return models.BooleanNode{Value: false}
	case "NIL":
		return models.NilNode{}
	case "STRING":
		joinedString := strings.Join(splitToken, " ")
		return models.StringNode{Value: strings.Split(joinedString, "\"")[1]}
	default:
		return nil
	}
}

func parseParrenthesisExpr(tokens []string) models.Node {
	innerTokens := tokens[1 : len(tokens)-1]
	var innerNode models.Node
	if len(innerTokens) == 1 {
		splitToken := strings.Split(innerTokens[0], " ")
		innerNode = parsevalue(splitToken)
	} else if utils.IsUraryExpr(innerTokens) {
		innerNode = parseUraryExpr(innerTokens)
	} else if utils.IsBinaryExpression(innerTokens) {
		innerNode = parseBinaryExpr(innerTokens)
	} else if utils.IsParethesizedExpr(innerTokens) {
		innerNode = parseParrenthesisExpr(innerTokens)
	} else {
		innerNode = models.NilNode{}
	}

	result := "(group " + innerNode.String() + ")"
	return models.StringNode{
		Value: result,
	}
}

func parseUraryExpr(tokens []string) models.Node {
	splitToken := strings.Split(tokens[0], " ")
	operator := splitToken[1]

	var operand models.Node
	remainingTokens := tokens[1:]
	if utils.IsUraryExpr(remainingTokens) {
		operand = parseUraryExpr(remainingTokens)
	} else if utils.IsParethesizedExpr(remainingTokens) {
		operand = parseParrenthesisExpr(remainingTokens)
	} else if len(remainingTokens) == 1 {
		splitRemain0 := strings.Split(remainingTokens[0], " ")
		operand = parsevalue(splitRemain0)
	} else {
		return models.NilNode{}
	}
	if operand.Evaluate() == nil || operand.String() == "<nil>" {
		return models.NilNode{}
	}
	result := "(" + operator + " " + operand.String() + ")"
	return models.StringNode{
		Value: result,
	}
}
