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
	if utils.IsUnaryExpr(tokens.Success) {
		return parseUnaryExpr(tokens.Success)
	}
	return models.NilNode{}
}

func parseBinaryExpr(tokens []string) models.Node {
	if utils.IsSingleBinary(tokens) {
		return parseSingleBinaryExpr(tokens)
	} else {
		return parseMultipleBinaryExpr(tokens)
	}
}

func parseSingleBinaryExpr(tokens []string) models.Node {
	var left, right models.Node
	op := ""
	if utils.Isoperator(tokens[0]) {
		left = parseUnaryExpr(tokens[:2])
		splitOperator := strings.Split(tokens[2], " ")
		op = parseOperator(splitOperator)
		if len(tokens[3:]) == 1 {
			splitValue := strings.Split(tokens[3], " ")
			right = parsevalue(splitValue)
		} else {
			right = parseUnaryExpr(tokens[3:])
		}
	} else {
		splitValue := strings.Split(tokens[0], " ")
		left = parsevalue(splitValue)
		splitOperator := strings.Split(tokens[1], " ")
		op = parseOperator(splitOperator)
		if len(tokens[2:]) == 1 {
			splitValue = strings.Split(tokens[2], " ")
			right = parsevalue(splitValue)
		} else {
			right = parseUnaryExpr(tokens[2:])
		}
	}
	result := models.BinaryNode{
		Left:  left,
		Op:    op,
		Right: right,
	}

	return result
}

func parseMultipleBinaryExpr(tokens []string) models.Node {
	precedence := map[string]int{
		"*": 4,
		"/": 3,
		"+": 2,
		"-": 1,
	}
	var left models.Node
	var right models.Node
	currentPosition := 0
	if utils.Isoperator(tokens[0]) {
		left = parseUnaryExpr(tokens[:2])
		currentPosition = 2
	} else {
		splitValue := strings.Split(tokens[0], " ")
		left = parsevalue(splitValue)
		currentPosition = 1
	}
	splitOperator := strings.Split(tokens[currentPosition], " ")
	op := parseOperator(splitOperator)
	currentPosition++
	if utils.Isoperator(tokens[currentPosition]) {
		right = parseUnaryExpr(tokens[currentPosition:2])
		currentPosition += 2
	} else {
		splitValue := strings.Split(tokens[currentPosition], " ")
		right = parsevalue(splitValue)
		currentPosition += 1
	}
	previousBinary := models.BinaryNode{
		Left:  left,
		Op:    op,
		Right: right,
	}
	for currentPosition < len(tokens) {
		splitOperator = strings.Split(tokens[currentPosition], " ")
		op = parseOperator(splitOperator)
		currentPosition++
		if currentPosition >= len(tokens) {
			return models.NilNode{}
		}
		if utils.Isoperator(tokens[currentPosition]) {
			right = parseUnaryExpr(tokens[currentPosition : currentPosition+2])
			currentPosition++
		} else {
			if strings.HasPrefix(tokens[currentPosition], "LEFT_PAREN") {
				var parenEnd = 0
				for i := currentPosition; i < len(tokens); i++ {
					if strings.HasPrefix(tokens[i], "RIGHT_PAREN") {
						parenEnd = i
						break
					}
				}
				if parenEnd == 0 {
					return models.NilNode{}
				}
				right = parseParrenthesisExpr(tokens[currentPosition : parenEnd+1])
				currentPosition = parenEnd
			} else {
				splitValue := strings.Split(tokens[currentPosition], " ")
				right = parsevalue(splitValue)
			}
		}
		if precedence[op] > precedence[previousBinary.Op] {
			temp := previousBinary
			previousBinary = models.BinaryNode{
				Left:  previousBinary.Right,
				Op:    op,
				Right: right,
			}
			left = models.ShiftedBinaryNode{
				Value: models.BinaryNode{
					Left:  previousBinary,
					Op:    temp.Op,
					Right: temp.Left,
				},
			}
			//todo to ensure consitency btwn the left and previousBinary
			previousBinary = models.BinaryNode{
				Left:  previousBinary,
				Op:    temp.Op,
				Right: temp.Left,
			}
		} else {
			left = models.BinaryNode{
				Left:  previousBinary,
				Op:    op,
				Right: right,
			}
			previousBinary = models.BinaryNode{
				Left:  previousBinary,
				Op:    op,
				Right: right,
			}
		}
		currentPosition++

	}
	return left
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
		return models.NilNode{}
	}
}

func parseParrenthesisExpr(tokens []string) models.Node {
	innerTokens := tokens[1 : len(tokens)-1]
	var innerNode models.Node
	if len(innerTokens) == 1 {
		splitToken := strings.Split(innerTokens[0], " ")
		innerNode = parsevalue(splitToken)
	} else if utils.IsUnaryExpr(innerTokens) {
		innerNode = parseUnaryExpr(innerTokens)
	} else if utils.IsBinaryExpression(innerTokens) {
		innerNode = parseBinaryExpr(innerTokens)
	} else if utils.IsParethesizedExpr(innerTokens) {
		innerNode = parseParrenthesisExpr(innerTokens)
	} else {
		innerNode = models.NilNode{}
	}

	return models.ParenthesisNode{
		Expression: innerNode,
	}
}

func parseUnaryExpr(tokens []string) models.Node {
	splitToken := strings.Split(tokens[0], " ")
	operator := splitToken[1]

	var operand models.Node
	if len(tokens[1:]) == 1 {
		splitRemain0 := strings.Split(tokens[1], " ")
		operand = parsevalue(splitRemain0)
	} else {
		remainingTokens := tokens[1:]
		if utils.IsUnaryExpr(remainingTokens) {
			operand = parseUnaryExpr(remainingTokens)
		} else if utils.IsParethesizedExpr(remainingTokens) {
			operand = parseParrenthesisExpr(remainingTokens)
		} else {
			return models.NilNode{}
		}
	}
	if operand.Evaluate() == nil || operand.String() == "<nil>" {
		return models.NilNode{}
	}
	return models.UnaryNode{
		Op:    operator,
		Value: operand,
	}
}
