package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Stan-breaks/app/models"
)

func Parse(tokens models.Tokens) models.Node {
	var left, right models.Node
	op := ""
	numToken := len(tokens.Success)
	if numToken == 3 && strings.Split(tokens.Success[0], " ")[0] == "NUMBER" {
		for index, item := range tokens.Success {
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
	} else {
		if numToken == 1 {
			splitToken := strings.Split(tokens.Success[0], " ")
			return parsevalue(splitToken)
		}
		return parseParrenthesisExpr(tokens.Success)
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
			paren = 1
		case "RIGHT_PAREN":
			if paren == 0 {
				return models.NilNode{}
			}
		case "STRING", "NUMBER", "TRUE", "FALSE":
			value = parsevalue(splitToken)
		default:
			return models.NilNode{}
		}
	}
	return models.StringNode{
		Value: "(group " + value.String() + ")",
	}

}
