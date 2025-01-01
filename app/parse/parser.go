package parse

import (
	"fmt"
	"sort"
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
	var left models.Node
	currentPosition := 0
	myTokens := rearrangeBinary(tokens)
	fmt.Print(myTokens)
	if utils.Isoperator(myTokens[0]) {
		left = parseUnaryExpr(myTokens[:2])
		currentPosition = 2
	} else {
		splitValue := strings.Split(myTokens[0], " ")
		left = parsevalue(splitValue)
		currentPosition = 1
	}
	for currentPosition < len(myTokens) {
		splitOperator := strings.Split(myTokens[currentPosition], " ")
		op := parseOperator(splitOperator)
		currentPosition++
		var right models.Node
		if currentPosition >= len(myTokens) {
			return models.NilNode{}
		}
		if utils.Isoperator(myTokens[currentPosition]) {
			right = parseUnaryExpr(myTokens[currentPosition : currentPosition+2])
			currentPosition++
		} else {
			if strings.HasPrefix(myTokens[currentPosition], "LEFT_PAREN") {
				var parenEnd = 0
				for i := currentPosition; i < len(myTokens); i++ {
					if strings.HasPrefix(myTokens[i], "RIGHT_PAREN") {
						parenEnd = i
						break
					}
				}
				if parenEnd == 0 {
					return models.NilNode{}
				}
				right = parseParrenthesisExpr(myTokens[currentPosition : parenEnd+1])
				currentPosition = parenEnd
			} else {
				splitValue := strings.Split(myTokens[currentPosition], " ")
				right = parsevalue(splitValue)
			}
		}
		currentPosition++
		left = models.BinaryNode{
			Left:  left,
			Op:    op,
			Right: right,
		}
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

func rearrangeBinary(tokens []string) []string {
	precedence := map[string]int{
		"STAR":  4, // *
		"SLASH": 3, // /
		"PLUS":  2, // +
		"MINUS": 1, // -
	}
	var operators []models.Operator
	currentPosition := 0
	if utils.Isoperator(tokens[0]) {
		currentPosition = 2
	} else {
		currentPosition = 1
	}
	for currentPosition < len(tokens) {
		splitToken := strings.Split(tokens[currentPosition], " ")
		if prec, ok := precedence[splitToken[0]]; ok {
			operators = append(operators, models.Operator{
				Index:      currentPosition,
				Precedence: prec,
				Token:      tokens[currentPosition],
			})
		}
		if currentPosition < len(tokens)-1 && utils.Isoperator(tokens[currentPosition+1]) {
			currentPosition += 1
		}
		currentPosition += 2
	}

	if len(operators) == 0 {
		return tokens
	}

	sort.SliceStable(operators, func(i, j int) bool {
		return operators[i].Precedence > operators[j].Precedence
	})

	result := make([]string, len(tokens))
	copy(result, tokens)
	for _, op := range operators {
		leftStart, leftEnd := findLeftOperand(result, op.Index)
		rightStart, rightEnd := findRightOperand(result, op.Index)

		// Create new arrangement
		newArrangement := make([]string, 0)

		// Add tokens before this expression
		newArrangement = append(newArrangement, result[:leftStart]...)

		// Add this expression
		newArrangement = append(newArrangement, result[leftStart:leftEnd]...)
		newArrangement = append(newArrangement, op.Token)
		newArrangement = append(newArrangement, result[rightStart:rightEnd]...)

		// Add remaining tokens
		newArrangement = append(newArrangement, result[rightEnd:]...)

		result = newArrangement
	}

	return result
}

// Helper function to find the left operand boundaries
func findLeftOperand(tokens []string, opIndex int) (start, end int) {
	end = opIndex
	start = end - 1

	// Handle parenthesized expressions
	if strings.HasPrefix(tokens[start], "RIGHT_PAREN") {
		parenCount := 1
		for start > 0 {
			start--
			if strings.HasPrefix(tokens[start], "RIGHT_PAREN") {
				parenCount++
			} else if strings.HasPrefix(tokens[start], "LEFT_PAREN") {
				parenCount--
				if parenCount == 0 {
					break
				}
			}
		}
	}

	return start, end
}

// Helper function to find the right operand boundaries
func findRightOperand(tokens []string, opIndex int) (start, end int) {
	start = opIndex + 1
	end = start + 1

	// Handle parenthesized expressions
	if strings.HasPrefix(tokens[start], "LEFT_PAREN") {
		parenCount := 1
		for end < len(tokens) {
			if strings.HasPrefix(tokens[end], "LEFT_PAREN") {
				parenCount++
			} else if strings.HasPrefix(tokens[end], "RIGHT_PAREN") {
				parenCount--
				if parenCount == 0 {
					end++
					break
				}
			}
			end++
		}
	}

	return start, end
}
