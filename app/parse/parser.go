package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Stan-breaks/app/environment"
	"github.com/Stan-breaks/app/models"
	nativefunctions "github.com/Stan-breaks/app/nativeFunctions"
	"github.com/Stan-breaks/app/utils"
)

func Parse(tokens []models.TokenInfo) (models.Node, []string) {
	if len(tokens) == 0 {
		return models.StringNode{Value: ""}, nil
	}
	if len(tokens) == 1 {
		return parsevalue(tokens[0])
	}
	if utils.IsParethesizedExpr(tokens) {
		return parseParrenthesisExpr(tokens)
	}
	if utils.IsUnaryExpr(tokens) {
		return parseUnaryExpr(tokens)
	}
	if utils.IsBinaryExpression(tokens) {
		return parseBinaryExpr(tokens)
	}
	var arrErr []string
	splitToken := strings.Split(tokens[0].Token, " ")
	errstr := fmt.Sprintf("[line %d] Error at %s", tokens[0].Line, splitToken[1])
	arrErr = append(arrErr, errstr)
	return models.StringNode{Value: ""}, arrErr
}

func parseBinaryExpr(tokens []models.TokenInfo) (models.Node, []string) {
	if utils.IsSingleBinary(tokens) {
		return parseSingleBinaryExpr(tokens)
	} else {
		return parseMultipleBinaryExpr(tokens)
	}
}

func parseSingleBinaryExpr(tokens []models.TokenInfo) (models.Node, []string) {
	var left models.Node
	var right models.Node
	currentPosition := 0
	var arrErr []string
	var err []string
	if utils.Isoperator(tokens[currentPosition]) {
		left, err = parseUnaryExpr(tokens[currentPosition : currentPosition+2])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition = 2
	} else if strings.HasPrefix(tokens[currentPosition].Token, "LEFT_PAREN") {
		var parenEnd = 0
		for i := currentPosition; i < len(tokens); i++ {
			if strings.HasPrefix(tokens[i].Token, "RIGHT_PAREN") {
				parenEnd = i
				break
			}
		}
		if parenEnd == 0 {
			splitToken := strings.Split(tokens[currentPosition].Token, " ")
			errstr := fmt.Sprintf("[line %d] Error at %s", tokens[currentPosition].Line, splitToken[1])
			arrErr = append(arrErr, errstr)
			return models.StringNode{Value: ""}, arrErr

		}
		left, err = parseParrenthesisExpr(tokens[currentPosition : parenEnd+1])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition = parenEnd + 1
	} else {
		left, err = parsevalue(tokens[0])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition = 1
	}
	splitOperator := strings.Split(tokens[currentPosition].Token, " ")
	op := parseOperator(splitOperator)
	currentPosition++
	if utils.Isoperator(tokens[currentPosition]) {
		right, err = parseUnaryExpr(tokens[currentPosition:])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
	} else if strings.HasPrefix(tokens[currentPosition].Token, "LEFT_PAREN") {
		var parenEnd = 0
		for i := currentPosition; i < len(tokens); i++ {
			if strings.HasPrefix(tokens[i].Token, "RIGHT_PAREN") {
				parenEnd = i
				break
			}
		}
		if parenEnd == 0 {
			splitToken := strings.Split(tokens[currentPosition].Token, " ")
			errstr := fmt.Sprintf("[line %d] Error at %s", tokens[currentPosition].Line, splitToken[1])
			arrErr = append(arrErr, errstr)
			return models.StringNode{Value: ""}, arrErr
		}
		right, err = parseParrenthesisExpr(tokens[currentPosition:])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
	} else {
		right, err = parsevalue(tokens[currentPosition])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
	}
	result := models.BinaryNode{
		Left:  left,
		Op:    op,
		Right: right,
	}
	if len(arrErr) == 0 {
		return result, nil
	} else {
		return result, arrErr
	}
}

func parseMultipleBinaryExpr(tokens []models.TokenInfo) (models.Node, []string) {
	precedence := map[string]int{
		"*": 4,
		"/": 3,
		"+": 2,
		"-": 1,
	}
	var left models.Node
	var right models.Node
	var arrErr []string
	var err []string
	currentPosition := 0
	if utils.Isoperator(tokens[0]) {
		left, err = parseUnaryExpr(tokens[:2])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition = 2
	} else if strings.HasPrefix(tokens[0].Token, "LEFT_PAREN") {
		var parenEnd = 0
		for i := currentPosition; i < len(tokens); i++ {
			if strings.HasPrefix(tokens[i].Token, "RIGHT_PAREN") {
				parenEnd = i
				break
			}
		}
		if parenEnd == 0 {
			splitToken := strings.Split(tokens[currentPosition].Token, " ")
			errstr := fmt.Sprintf("[line %d] Error at %s", tokens[currentPosition].Line, splitToken[1])
			arrErr = append(arrErr, errstr)
			return models.StringNode{Value: ""}, arrErr
		}
		left, err = parseParrenthesisExpr(tokens[currentPosition : parenEnd+1])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition = parenEnd + 1
	} else {
		left, err = parsevalue(tokens[0])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition = 1
	}
	splitOperator := strings.Split(tokens[currentPosition].Token, " ")
	op := parseOperator(splitOperator)
	currentPosition++
	if utils.Isoperator(tokens[currentPosition]) {
		right, err = parseUnaryExpr(tokens[currentPosition : currentPosition+2])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition += 2
	} else if strings.HasPrefix(tokens[currentPosition].Token, "LEFT_PAREN") {
		var parenEnd = 0
		for i := currentPosition; i < len(tokens); i++ {
			if strings.HasPrefix(tokens[i].Token, "RIGHT_PAREN") {
				parenEnd = i
				break
			}
		}
		if parenEnd == 0 {
			splitToken := strings.Split(tokens[currentPosition].Token, " ")
			errstr := fmt.Sprintf("[line %d] Error at %s", tokens[currentPosition].Line, splitToken[1])
			arrErr = append(arrErr, errstr)
			return models.StringNode{Value: ""}, arrErr
		}
		right, err = parseParrenthesisExpr(tokens[currentPosition : parenEnd+1])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition = parenEnd + 1
	} else {
		right, err = parsevalue(tokens[currentPosition])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		currentPosition += 1
	}
	previousBinary := models.BinaryNode{
		Left:  left,
		Op:    op,
		Right: right,
	}
	for currentPosition < len(tokens) {
		splitOperator = strings.Split(tokens[currentPosition].Token, " ")
		op = parseOperator(splitOperator)
		currentPosition++
		if currentPosition >= len(tokens) {
			splitToken := strings.Split(tokens[currentPosition].Token, " ")
			errstr := fmt.Sprintf("[line %d] Error at %s", tokens[currentPosition].Line, splitToken[1])
			arrErr = append(arrErr, errstr)
			return models.NilNode{}, arrErr
		}
		if utils.Isoperator(tokens[currentPosition]) {
			right, err = parseUnaryExpr(tokens[currentPosition : currentPosition+2])
			if err != nil {
				arrErr = append(arrErr, err...)
			}
			currentPosition++
		} else {
			if strings.HasPrefix(tokens[currentPosition].Token, "LEFT_PAREN") {
				var parenEnd = 0
				for i := currentPosition; i < len(tokens); i++ {
					if strings.HasPrefix(tokens[i].Token, "RIGHT_PAREN") {
						parenEnd = i
						break
					}
				}
				if parenEnd == 0 {
					splitToken := strings.Split(tokens[currentPosition].Token, " ")
					errstr := fmt.Sprintf("[line %d] Error at %s", tokens[currentPosition].Line, splitToken[1])
					arrErr = append(arrErr, errstr)
					return models.StringNode{Value: ""}, arrErr
				}
				right, err = parseParrenthesisExpr(tokens[currentPosition : parenEnd+1])
				if err != nil {
					arrErr = append(arrErr, err...)
				}
				currentPosition = parenEnd
			} else {
				right, err = parsevalue(tokens[currentPosition])
				if err != nil {
					arrErr = append(arrErr, err...)
				}
			}
		}
		if precedence[op] > precedence[previousBinary.Op] {
			temp := previousBinary
			previousBinary = models.BinaryNode{
				Left:  previousBinary.Right,
				Op:    op,
				Right: right,
			}
			left = models.BinaryNode{
				Left:    previousBinary,
				Op:      temp.Op,
				Right:   temp.Left,
				Shifted: 1,
			}
			previousBinary = left.(models.BinaryNode)
		} else {
			left = models.BinaryNode{
				Left:  previousBinary,
				Op:    op,
				Right: right,
			}
			previousBinary = left.(models.BinaryNode)
		}
		currentPosition++
	}
	if len(arrErr) == 0 {
		return left, nil
	} else {
		return left, arrErr
	}
}

func parseOperator(splitToken []string) string {
	switch splitToken[0] {
	case "PLUS", "MINUS", "STAR", "SLASH", "BANG_EQUAL", "EQUAL_EQUAL",
		"LESS", "GREATER", "LESS_EQUAL", "GREATER_EQUAL", "AND", "OR":
		return splitToken[1]
	default:
		return ""
	}
}

func parsevalue(token models.TokenInfo) (models.Node, []string) {
	splitToken := strings.Split(token.Token, " ")
	switch splitToken[0] {
	case "NUMBER":
		num, _ := strconv.ParseFloat(splitToken[1], 32)
		floatnum := float32(num)
		return models.NumberNode{Value: floatnum}, nil
	case "TRUE":
		return models.BooleanNode{Value: true}, nil
	case "FALSE":
		return models.BooleanNode{Value: false}, nil
	case "NIL":
		return models.NilNode{}, nil
	case "STRING":
		joinedString := strings.Join(splitToken, " ")
		return models.StringNode{Value: strings.Split(joinedString, "\"")[1]}, nil
	case "IDENTIFIER":
		valname := splitToken[1]
		value, exist := environment.Global.Get(valname)
		if !exist {
			err := fmt.Sprintf("[Line %d] Error at %s", token.Line, splitToken[1])
			var errors []string
			errors = append(errors, err)
			return models.NilNode{}, errors
		}
		switch v := value.(type) {
		case bool:
			return models.BooleanNode{Value: v}, nil
		case string:
			return models.StringNode{Value: v}, nil
		case float32:
			return models.NumberNode{Value: v}, nil
		case int:
			return models.NumberNode{Value: float32(v)}, nil
		default:
			err := fmt.Sprintf("[Line %d] Error at %s", token.Line, splitToken[1])
			var errors []string
			errors = append(errors, err)
			return models.NilNode{}, errors
		}
	case "FUNCTION":
		var errors []string
		funcName := splitToken[1]
		value := nativefunctions.GlobalFunctions[funcName]
		return models.NumberNode{Value: value.(float32)}, errors
	default:
		err := fmt.Sprintf("[Line %d] Error at %s", token.Line, splitToken[1])
		var errors []string
		errors = append(errors, err)
		return models.StringNode{Value: ""}, errors
	}
}

func parseParrenthesisExpr(tokens []models.TokenInfo) (models.Node, []string) {
	innerTokens := tokens[1 : len(tokens)-1]
	var innerNode models.Node
	var arrErr []string
	var err []string
	if len(innerTokens) == 1 {
		innerNode, err = parsevalue(innerTokens[0])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
	} else if utils.IsBinaryExpression(innerTokens) {
		innerNode, err = parseBinaryExpr(innerTokens)
		if err != nil {
			arrErr = append(arrErr, err...)
		}
	} else if utils.IsUnaryExpr(innerTokens) {
		innerNode, err = parseUnaryExpr(innerTokens)
		if err != nil {
			arrErr = append(arrErr, err...)
		}
	} else if utils.IsParethesizedExpr(innerTokens) {
		innerNode, err = parseParrenthesisExpr(innerTokens)
		if err != nil {
			arrErr = append(arrErr, err...)
		}
	} else {
		splitToken := strings.Split(tokens[len(tokens)-1].Token, " ")
		err := fmt.Sprintf("[Line %d] Error at '%s': Expect expression.", innerTokens[0].Line, splitToken[1])
		var errors []string
		arrErr = append(errors, err)
		return models.StringNode{Value: ""}, arrErr
	}
	return models.ParenthesisNode{
		Expression: innerNode,
	}, arrErr
}

func parseUnaryExpr(tokens []models.TokenInfo) (models.Node, []string) {
	splitToken := strings.Split(tokens[0].Token, " ")
	operator := splitToken[1]
	var operand models.Node
	var arrErr []string
	var err []string
	if len(tokens[1:]) == 1 {
		operand, err = parsevalue(tokens[1])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
	} else {
		remainingTokens := tokens[1:]
		if utils.IsParethesizedExpr(remainingTokens) {
			operand, err = parseParrenthesisExpr(remainingTokens)
			if err != nil {
				arrErr = append(arrErr, err...)
			}
		} else if utils.IsUnaryExpr(remainingTokens) {
			operand, err = parseUnaryExpr(remainingTokens)
			if err != nil {
				arrErr = append(arrErr, err...)
			}
		} else {
			splitToken = strings.Split(tokens[1].Token, " ")
			errstr := fmt.Sprintf("[line %d] Error at %s", tokens[1].Line, splitToken[1])
			arrErr = append(arrErr, errstr)
			return models.StringNode{Value: ""}, arrErr
		}
	}
	if operand.Evaluate() == nil || operand.String() == "<nil>" {
		splitToken = strings.Split(tokens[1].Token, " ")
		errstr := fmt.Sprintf("[line %d] Error at %s", tokens[1].Line, splitToken[1])
		arrErr = append(arrErr, errstr)
		return models.StringNode{Value: ""}, arrErr
	}
	return models.UnaryNode{
		Op:    operator,
		Value: operand,
	}, nil
}
