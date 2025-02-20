package parse

import (
	"fmt"

	"github.com/Stan-breaks/app/environment"
	"github.com/Stan-breaks/app/models"
	"github.com/Stan-breaks/app/nativeFunctions"
	"github.com/Stan-breaks/app/utils"
)

func Parse(tokens []models.Token) (models.Node, []string) {
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
	errstr := fmt.Sprintf("[line %d] Error at %s", tokens[0].Line, tokens[0].Lexem)
	arrErr = append(arrErr, errstr)
	return models.StringNode{Value: ""}, arrErr
}

func parseBinaryExpr(tokens []models.Token) (models.Node, []string) {
	if utils.IsSingleBinary(tokens) {
		return parseSingleBinaryExpr(tokens)
	}
	return parseMultipleBinaryExpr(tokens)

}

func parseSingleBinaryExpr(tokens []models.Token) (models.Node, []string) {
	currentPosition := 0
	var arrErr []string
	var err []string

	left, tokensUsed, err := parseOperand(tokens[currentPosition:])
	if len(err) == 0 {
		arrErr = append(arrErr, err...)
	}
	currentPosition += tokensUsed

	op := parseOperator(tokens[currentPosition])
	currentPosition++

	right, _, err := parseOperand(tokens[currentPosition:])
	if len(err) == 0 {
		arrErr = append(arrErr, err...)
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

func parseOperand(tokens []models.Token) (models.Node, int, []string) {
	var node models.Node
	var arrErr []string
	var err []string
	tokensUsed := 0
	if len(tokens) == 0 {
		arrErr = append(arrErr, "no tokens in operand")
		return models.NilNode{}, 0, arrErr
	}
	if utils.Isoperator(tokens[0]) {
		if tokens[1].Type == models.LEFT_PAREN {
			parenEnd := utils.FindClosingParen(tokens)
			tokensUsed = parenEnd + 1
		} else {
			tokensUsed = 2
		}
		node, err = parseUnaryExpr(tokens[:tokensUsed])
		if err != nil {
			arrErr = append(arrErr, err...)
		}

	} else if tokens[0].Type == models.LEFT_PAREN {
		parenEnd := utils.FindClosingParen(tokens)
		if parenEnd == 0 {
			errstr := fmt.Sprintf("[line %d] Error at %s", tokens[tokensUsed].Line, tokens[tokensUsed].Lexem)
			arrErr = append(arrErr, errstr)
			return models.StringNode{Value: ""}, tokensUsed, arrErr
		}
		node, err = parseParrenthesisExpr(tokens[tokensUsed : parenEnd+1])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		tokensUsed = parenEnd + 1
	} else {
		node, err = parsevalue(tokens[0])
		if err != nil {
			arrErr = append(arrErr, err...)
		}
		tokensUsed = 1
	}
	return node, tokensUsed, arrErr
}

func parseMultipleBinaryExpr(tokens []models.Token) (models.Node, []string) {
	precedence := map[string]int{
		"*":  4,
		"/":  4,
		"+":  3,
		"-":  3,
		">":  2,
		"<":  2,
		">=": 2,
		"<=": 2,
	}
	var nodeStack []models.BinaryNode
	var opStack []string
	var arrErr []string
	currentPosition := 0

	left, tokensUsed, err := parseOperand(tokens[currentPosition:])
	arrErr = append(arrErr, err...)
	currentPosition += tokensUsed

	for currentPosition < len(tokens) {
		currentOp := parseOperator(tokens[currentPosition])
		currentPosition++

		right, tokensUsed, err := parseOperand(tokens[currentPosition:])
		arrErr = append(arrErr, err...)
		currentPosition += tokensUsed

		for len(opStack) > 0 && precedence[currentOp] <= precedence[opStack[len(opStack)-1]] {
			prevOp := opStack[len(opStack)-1]
			opStack = opStack[:len(opStack)-1]

			prevRight := nodeStack[len(nodeStack)-1]
			nodeStack = nodeStack[:len(nodeStack)-1]

			left = models.BinaryNode{
				Left:  prevRight.Left,
				Op:    prevOp,
				Right: left,
			}
		}

		nodeStack = append(nodeStack, models.BinaryNode{Left: left, Op: currentOp})
		opStack = append(opStack, currentOp)
		left = right
	}

	for len(opStack) > 0 {
		prevOp := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]

		prevNode := nodeStack[len(nodeStack)-1]
		nodeStack = nodeStack[:len(nodeStack)-1]

		left = models.BinaryNode{
			Left:  prevNode.Left,
			Op:    prevOp,
			Right: left,
		}
	}

	if len(arrErr) > 0 {
		return left, arrErr
	}
	return left, nil
}

func parseOperator(token models.Token) string {
	switch token.Type {
	case models.PLUS, models.MINUS, models.STAR, models.SLASH, models.BANG_EQUAL, models.EQUAL_EQUAL,
		models.LESS, models.GREATER_EQUAL, models.LESS_EQUAL, models.GREATER, models.AND, models.OR:
		return token.Lexem
	default:
		return ""
	}
}

func parsevalue(token models.Token) (models.Node, []string) {
	switch token.Type {
	case models.NUMBER:
		var value float64
		switch v := token.Literal.(type) {
		case int:
			value = float64(v)
		case float64:
			value = v
		case float32:
			value = float64(v)
		}
		return models.NumberNode{Value: value}, nil
	case models.TRUE:
		return models.BooleanNode{Value: true}, nil
	case models.FALSE:
		return models.BooleanNode{Value: false}, nil
	case models.NIL:
		return models.NilNode{}, nil
	case models.STRING:
		return models.StringNode{Value: token.Literal.(string)}, nil
	case models.IDENTIFIER:
		valname := token.Lexem
		value, exist := environment.Global.Get(valname)
		if !exist {
			err := fmt.Sprintf("[Line %d] Undefined variable: %s", token.Line, valname)
			var errors []string
			errors = append(errors, err)
			return models.NilNode{}, errors
		}
		switch v := value.(type) {
		case bool:
			return models.BooleanNode{Value: v}, nil
		case string:
			return models.StringNode{Value: v}, nil
		case float64:
			return models.NumberNode{Value: v}, nil
		case int:
			return models.NumberNode{Value: float64(v)}, nil
		case models.Function:
			return models.StringNode{Value: "<fn " + valname + ">"}, nil
		default:
			err := fmt.Sprintf("[Line %d] Invalid token: %s", token.Line, valname)
			var errors []string
			errors = append(errors, err)
			return models.NilNode{}, errors
		}
	case "FUNCTION":
		var errors []string
		funcName := token.Lexem
		val := nativeFunctions.GlobalFunctions[funcName]
		fn := val.(func() float64)
		value := fn()
		return models.NumberNode{Value: value}, errors
	default:
		err := fmt.Sprintf("[Line %d] Error at %s", token.Line, token.Lexem)
		var errors []string
		errors = append(errors, err)
		return models.StringNode{Value: ""}, errors
	}
}

func parseParrenthesisExpr(tokens []models.Token) (models.Node, []string) {
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
		err := fmt.Sprintf("[Line %d] Error at '%s': Expect expression.", innerTokens[0].Line, innerTokens[0].Lexem)
		var errors []string
		arrErr = append(errors, err)
		return models.StringNode{Value: ""}, arrErr
	}
	return models.ParenthesisNode{
		Expression: innerNode,
	}, arrErr
}

func parseUnaryExpr(tokens []models.Token) (models.Node, []string) {
	operator := tokens[0].Lexem
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
			errstr := fmt.Sprintf("[line %d] Error at %s", tokens[1].Line, tokens[1].Lexem)
			arrErr = append(arrErr, errstr)
			return models.StringNode{Value: ""}, arrErr
		}
	}
	if operand.Evaluate() == nil || operand.String() == "<nil>" {
		errstr := fmt.Sprintf("[line %d] Error at %s", tokens[1].Line, tokens[1].Lexem)
		arrErr = append(arrErr, errstr)
		return models.StringNode{Value: ""}, arrErr
	}
	return models.UnaryNode{
		Op:    operator,
		Value: operand,
	}, nil
}
