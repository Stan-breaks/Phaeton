package parse

import (
	"strconv"
	"strings"

	"github.com/Stan-breaks/app/models"
)

func Parse(tokens models.Tokens) models.Node {
	var left, right models.Node
	op := ""
	numToken := len(tokens.Success)
	if numToken == 3 {
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

	}
	if numToken == 1 {
		splitToken := strings.Split(tokens.Success[0], " ")
		return parsevalue(splitToken)
	}

	return nil
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
	default:
		return nil
	}
}
