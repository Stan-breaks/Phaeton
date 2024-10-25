package parse

import (
	"strconv"
	"strings"

	"github.com/Stan-breaks/app/models"
)

func Parse(tokens models.Tokens) models.Node {
	if len(tokens.Success) == 1 {
		splitToken := strings.Split(tokens.Success[0], " ")
		return parsevalue(splitToken)
	}
	if isBinaryExpression(tokens.Success) {
		return parseBinaryExpr(tokens.Success)
	}

	if node := parseUraryExpr(tokens.Success); node.Evaluate() != nil {
		return node
	}
	return parseParrenthesisExpr(tokens.Success)
}

func isBinaryExpression(tokens []string) bool {
	return len(tokens) == 3 && strings.HasPrefix(tokens[0], "NUMBER")
}

func parseBinaryExpr(tokens []string) models.Node {
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
		return models.NilNode{Value: "nil"}
	case "STRING":
		joinedString := strings.Join(splitToken, " ")
		return models.StringNode{Value: strings.Split(joinedString, "\"")[1]}
	default:
		return nil
	}
}

func parseParrenthesisExpr(tokens []string) models.Node {
	var value models.Node
	paren := 0
	for _, item := range tokens {
		splitToken := strings.Split(item, " ")
		switch splitToken[0] {
		case "LEFT_PAREN":
			paren += 1
		case "RIGHT_PAREN":
			if paren == 0 {
				return models.NilNode{}
			}
		case "STRING", "NUMBER", "TRUE", "FALSE", "NIL":
			value = parsevalue(splitToken)
		default:
			return models.NilNode{}
		}
	}
	result := ""
	for i := 0; i < paren-1; i++ {
		result += "(group "
	}
	result += ("(group " + value.String() + ")")
	for i := 0; i < paren-1; i++ {
		result += ")"
	}
	return models.StringNode{
		Value: result,
	}

}

func parseUraryExpr(tokens []string) models.Node {
	result := ""
	switch len(tokens) {
	case 2:
		splitToken0 := strings.Split(tokens[0], " ")
		splitToken1 := strings.Split(tokens[1], " ")
		value := parsevalue(splitToken1)
		result = "(" + splitToken0[1] + " " + value.String() + ")"
	case 3:
		splitToken0 := strings.Split(tokens[0], " ")
		splitToken1 := strings.Split(tokens[1], " ")
		splitToken2 := strings.Split(tokens[2], " ")
		value := parsevalue(splitToken2)
		result = "(" + splitToken0[1] + " (" + splitToken1[1] + " " + value.String() + "))"
	default:
		return models.NilNode{}
	}
	return models.StringNode{Value: result}
}
