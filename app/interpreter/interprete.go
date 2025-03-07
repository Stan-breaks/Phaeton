package interpreter

import (
	"fmt"
	"strings"

	"github.com/Stan-breaks/app/environment"
	"github.com/Stan-breaks/app/models"
	"github.com/Stan-breaks/app/parse"
	"github.com/Stan-breaks/app/tokenize"
	"github.com/Stan-breaks/app/utils"
)

func Interprete(tokens []models.Token) (interface{}, error) {
	currentPosition := 0
	for currentPosition < len(tokens) {
		if currentPosition >= len(tokens) {
			break
		}
		token := tokens[currentPosition]
		switch token.Type {
		case models.IF:
			tokensProcessed, ret, err := handleIf(tokens[currentPosition:])
			if err != nil {
				return nil, err
			}
			if ret != nil {
				return ret, nil
			}
			currentPosition += tokensProcessed
		case models.PRINT:
			tokensProcessed, err := handlePrint(tokens[currentPosition:])
			if err != nil {
				return nil, err
			}
			currentPosition += tokensProcessed
		case models.VAR:
			tokensProcessed, err := handleAssignment(tokens[currentPosition:])
			if err != nil {
				return nil, err
			}
			currentPosition += tokensProcessed
		case models.IDENTIFIER:
			if utils.IsFunctionCall(tokens[currentPosition:]) {
				_, tokensProcessed, err := handleFunCall(tokens[currentPosition:])
				if err != nil {
					return nil, err
				}
				currentPosition += tokensProcessed
			} else {
				tokensProcessed, err := handleReassignment(tokens[currentPosition:])
				if err != nil {
					return nil, err
				}
				currentPosition += tokensProcessed
			}

		case models.LEFT_PAREN:
			tokensProcessed, err := handleParenStatement(tokens[currentPosition:])
			if err != nil {
				return nil, err
			}
			currentPosition += tokensProcessed
		case models.WHILE:
			tokensProcessed, ret, err := handleWhile(tokens[currentPosition:])
			if err != nil {
				return nil, err
			}
			if ret != nil {
				return ret, nil
			}
			currentPosition += tokensProcessed
		case models.FOR:
			tokensProcessed, ret, err := handleFor(tokens[currentPosition:])
			if err != nil {
				return nil, err
			}
			if ret != nil {
				return ret, nil
			}
			currentPosition += tokensProcessed
		case models.FUN:
			tokensProcessed, err := handleFun(tokens[currentPosition:])
			if err != nil {
				return nil, err
			}
			currentPosition += tokensProcessed
		case models.RETURN:
			_, ret, err := handleReturn(tokens[currentPosition:])
			if err != nil {
				return nil, err
			}
			return ret, nil
		default:
			currentPosition++
		}
	}
	return nil, nil
}

func handleReturn(tokens []models.Token) (int, interface{}, error) {
	semicolon := utils.FindLastSemicolonInSameLine(tokens)
	if semicolon == -1 {
		return 0, nil, fmt.Errorf("no semicolon in return statement")
	}
	if len(tokens[1:semicolon]) == 0 {
		return semicolon + 1, nil, nil
	}
	result, err := handleExpression(tokens[1:semicolon])
	if err != nil {
		return 0, nil, fmt.Errorf("error with parsing return statement: %v", err.Error())
	}
	return semicolon + 1, result.Evaluate(), nil
}

func handleFunCall(tokens []models.Token) (models.Node, int, error) {
	environment.Global.PushScope()
	defer environment.Global.PopScope()
	funName := tokens[0].Lexem
	value, bool := environment.Global.Get(funName)
	if !bool {
		return models.NilNode{}, 0, fmt.Errorf("function not defined")
	}
	switch v := value.(type) {
	case models.Function:
		switch a := v.Arguments.(type) {
		case []models.Token:
			argumentEnd := utils.FindLastSemicolonInSameLine(tokens)
			if argumentEnd == -1 {
				argumentEnd = len(tokens) - 1
			} else {
				argumentEnd--
			}
			args := tokens[2:argumentEnd]
			arrA := utils.FindNoOfArgs(a)
			arrArgs := utils.FindNoOfArgs(args)
			if len(arrA) != len(arrArgs) {
				return models.NilNode{}, 0, fmt.Errorf("expected %d arguments, got %d", len(arrA), len(arrArgs))
			}
			if len(arrArgs) > 0 {
				for i := 0; i < len(arrArgs); i++ {
					valName := arrA[i][0].Lexem
					val, err := parse.Parse(arrArgs[i])
					if err != nil {
						return models.NilNode{}, 0, fmt.Errorf("invalid function arguments")
					}
					environment.Global.Set(valName, val.Evaluate())
				}
			}
		}
		switch b := v.Body.(type) {
		case []models.Token:
			ret, err := Interprete(b)
			if err != nil {
				return models.NilNode{}, 0, err
			}
			if ret != nil {
				switch v := ret.(type) {
				case float64:
					return models.NumberNode{Value: v}, len(tokens), nil
				case string:
					if strings.HasPrefix(v, "<fn") {
						funcName := strings.Trim(strings.Split(v, " ")[1], ">")
						value, bool := environment.Global.Get(funcName)
						if bool {
							outerScope := &environment.Global.Scopes[len(environment.Global.Scopes)-2]
							outerScope.Variables[funcName] = value
						}
					}
					return models.StringNode{Value: v}, len(tokens), nil
				}
			}
		}

	}
	return models.NilNode{}, len(tokens), nil
}

func handleFun(tokens []models.Token) (int, error) {
	positions := findFunPositions(tokens)
	if !positions.IsValid() {
		return positions.BodyEnd + 1, fmt.Errorf("invalid function")
	}
	funName := tokens[1].Lexem
	function := models.Function{
		Arguments: tokens[positions.ArgumentStart+1 : positions.ArgumentEnd],
		Body:      tokens[positions.BodyStart : positions.BodyEnd+1],
	}
	environment.Global.Set(funName, function)
	return positions.BodyEnd + 1, nil
}

func findFunPositions(tokens []models.Token) models.FunStatementPositions {
	positions := models.FunStatementPositions{
		ArgumentStart: -1,
		ArgumentEnd:   -1,
		BodyStart:     -1,
		BodyEnd:       -1,
	}
	parenCount := 0
	braceCount := 0
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Type {
		case models.LEFT_PAREN:
			if parenCount == 0 && positions.ArgumentStart == -1 {
				positions.ArgumentStart = i
			}
			parenCount++
		case models.RIGHT_PAREN:
			parenCount--
			if parenCount == 0 && positions.ArgumentEnd == -1 {
				positions.ArgumentEnd = i
				positions.BodyStart = i + 1

			}
		case models.LEFT_BRACE:
			braceCount++
		case models.RIGHT_BRACE:
			braceCount--
			if braceCount == 0 && positions.BodyEnd == -1 {
				positions.BodyEnd = i
			}
		}
	}
	return positions
}

func handleFor(tokens []models.Token) (int, interface{}, error) {
	environment.Global.PushScope()
	defer environment.Global.PopScope()

	positions := findForPositions(tokens)
	if !positions.IsValid() {
		return positions.BodyEnd + 1, nil, fmt.Errorf("invalid for statement")
	}
	declarationStart := positions.ConditionStart + 1
	declarationEnd := utils.FindSemicolonPosition(tokens[positions.ConditionStart+1:positions.ConditionEnd]) + positions.ConditionStart + 1
	conditionStart := declarationEnd + 1
	conditionEnd := utils.FindSemicolonPosition(tokens[declarationEnd+1:positions.ConditionEnd]) + declarationEnd + 1
	expressionStart := conditionEnd + 1
	expressionEnd := positions.ConditionEnd

	_, err := Interprete(tokens[declarationStart : declarationEnd+1])
	if err != nil {
		return positions.BodyEnd + 1, nil, err
	}
	condition, err := handleExpression(tokens[conditionStart:conditionEnd])
	if err != nil {
		return positions.BodyEnd + 1, nil, err
	}
	if condition.IsTruthy() {
		if expressionStart != expressionEnd {
			for {
				environment.Global.PushScope()
				ret, err := Interprete(tokens[positions.BodyStart : positions.BodyEnd+1])
				if err != nil {
					return positions.BodyEnd + 1, nil, err
				}
				environment.Global.PopScope()
				if ret != nil {
					return positions.BodyEnd + 1, ret, err
				}
				_, err = handleExpression(tokens[expressionStart:expressionEnd])
				if err != nil {
					return positions.BodyEnd + 1, nil, err
				}

				condition, err = handleExpression(tokens[conditionStart:conditionEnd])
				if err != nil {
					return positions.BodyEnd + 1, nil, err
				}
				if !condition.IsTruthy() {
					break
				}
			}
		} else {
			for {
				ret, err := Interprete(tokens[positions.BodyStart : positions.BodyEnd+1])
				if err != nil {
					return positions.BodyEnd + 1, nil, err
				}
				if ret != nil {
					return positions.BodyEnd + 1, ret, err
				}
				condition, err = handleExpression(tokens[conditionStart:conditionEnd])
				if err != nil {
					return positions.BodyEnd + 1, nil, err
				}
				if !condition.IsTruthy() {
					break
				}
			}
		}
	}
	return positions.BodyEnd + 1, nil, nil
}

func findForPositions(tokens []models.Token) models.ForStatementPositions {
	positions := models.ForStatementPositions{
		ConditionStart: 1,
		ConditionEnd:   -1,
		BodyStart:      -1,
		BodyEnd:        -1,
	}
	parenCount := 0
	braceCount := 0
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Type {
		case models.LEFT_PAREN:
			parenCount++
		case models.RIGHT_PAREN:
			parenCount--
			if parenCount == 0 && positions.ConditionEnd == -1 {
				positions.ConditionEnd = i
				positions.BodyStart = i + 1
				if tokens[i+1].Type != models.LEFT_BRACE {
					positions.BodyEnd = utils.FindLastSemicolonInSameLine(tokens[i+1:]) + i + 1
				}
			}
		case models.LEFT_BRACE:
			braceCount++
		case models.RIGHT_BRACE:
			braceCount--
			if braceCount == 0 && positions.BodyEnd == -1 {
				positions.BodyEnd = i
				goto exit
			}
		case models.SEMICOLON:
			if positions.BodyEnd == i && parenCount == 0 && braceCount == 0 {
				goto exit
			}
		}

	}
exit:
	return positions
}

func handleWhile(tokens []models.Token) (int, interface{}, error) {
	environment.Global.PushScope()
	defer environment.Global.PopScope()
	positions := findWhilePositions(tokens)
	if !positions.IsValid() {
		return positions.BodyEnd + 1, nil, fmt.Errorf("invalid while statement")
	}
	condition, err := handleExpression(tokens[positions.ConditionStart+1 : positions.ConditionEnd])
	if err != nil {
		return positions.BodyEnd + 1, nil, fmt.Errorf("invalid while condition")
	}
	if condition.IsTruthy() {
		for {
			ret, err := Interprete(tokens[positions.BodyStart : positions.BodyEnd+1])
			if err != nil {
				return positions.BodyEnd + 1, nil, fmt.Errorf("invalid while body")
			}
			if ret != nil {
				return positions.BodyEnd + 1, ret, err
			}
			condition, err = handleExpression(tokens[positions.ConditionStart+1 : positions.ConditionEnd])
			if err != nil {
				return positions.BodyEnd + 1, nil, fmt.Errorf("invalid while condition")
			}
			if !condition.IsTruthy() {
				break
			}
		}

	}
	return positions.BodyEnd + 1, nil, nil
}

func findWhilePositions(tokens []models.Token) models.WhileStatementPositions {
	positions := models.WhileStatementPositions{
		ConditionStart: 1,
		ConditionEnd:   -1,
		BodyStart:      -1,
		BodyEnd:        -1,
	}
	parenCount := 0
	braceCount := 0
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Type {
		case models.LEFT_PAREN:
			parenCount++
		case models.RIGHT_PAREN:
			parenCount--
			if parenCount == 0 && positions.ConditionEnd == -1 {
				positions.ConditionEnd = i
				positions.BodyStart = i + 1
				if tokens[i+1].Type != models.LEFT_BRACE {
					positions.BodyEnd = utils.FindLastSemicolonInSameLine(tokens[i+1:]) + i + 1
				}
			}
		case models.LEFT_BRACE:
			braceCount++
		case models.RIGHT_BRACE:
			braceCount--
			if braceCount == 0 && positions.BodyEnd == -1 {
				positions.BodyEnd = i
				goto exit
			}
		case models.SEMICOLON:
			if positions.BodyEnd == i && parenCount == 0 && braceCount == 0 {
				goto exit
			}
		}

	}
exit:
	return positions
}

func handleParenStatement(tokens []models.Token) (int, error) {
	end := utils.FindSemicolonPosition(tokens)
	if end == 0 {
		return 0, fmt.Errorf("no semicolon found")
	}
	startParen := 0
	endParen := utils.FindClosingParen(tokens[startParen:])
	for endParen < end {
		truthy, err := handleCondition(tokens[startParen+1 : endParen])
		if err != nil {
			return 0, err
		}
		op := endParen + 1
		if truthy && tokens[op].Type == models.OR {
			break
		} else if !truthy && tokens[op].Type == models.AND {
			break
		}
		startParen = endParen + 2
		endParen = utils.FindClosingParen(tokens[startParen:]) + startParen
	}
	return end, nil
}

func handleCondition(tokens []models.Token) (bool, error) {
	node, err := handleReassignmentCondition(tokens)
	if err != nil {
		return false, err
	}
	if node.IsTruthy() {
		return true, nil
	} else {
		return false, nil
	}
}

func handleAssignment(tokens []models.Token) (int, error) {
	if len(tokens) < 4 {
		return 0, fmt.Errorf("incomplete variable declaration")
	}
	nameToken := tokens[1]
	if nameToken.Type != models.IDENTIFIER {
		return 0, fmt.Errorf("no identifier found")
	}
	variableName := nameToken.Lexem
	if tokens[2].Type != models.EQUAL {
		return 0, fmt.Errorf("equal not found")
	}
	semicolonPosition := utils.FindSemicolonPosition(tokens[3:])
	if semicolonPosition == -1 {
		return 0, fmt.Errorf("no semicolon found")
	}
	if utils.IsFunctionCall(tokens[3 : semicolonPosition+4]) {
		result, _, err := handleFunCall(tokens[3 : semicolonPosition+3])
		if err != nil {
			return semicolonPosition + 4, err
		}
		funcName := strings.Trim(strings.Split(result.String(), " ")[1], ">")
		value, bool := environment.Global.Get(funcName)
		if bool {
			environment.Global.Set(variableName, value)
			currentScope := &environment.Global.Scopes[len(environment.Global.Scopes)-1]
			delete(currentScope.Variables, funcName)
		}
	} else {
		expression, err := parse.Parse(tokens[3 : semicolonPosition+3])
		if len(err) > 0 {
			return 0, fmt.Errorf("invalid assignment expression: %v", err[0])
		}
		value := expression.Evaluate()
		environment.Global.Set(variableName, value)
	}
	return semicolonPosition + 4, nil
}

func handleReassignment(tokens []models.Token) (int, error) {
	if tokens[1].Type != models.EQUAL {
		return 0, fmt.Errorf("no equal found in reassignment")
	}
	variableName := tokens[0].Lexem
	semicolonPosition := utils.FindSemicolonPosition(tokens[2:])
	if semicolonPosition == -1 {
		return 0, fmt.Errorf("no semicolon found in reassignment")
	}
	expression, err := parse.Parse(tokens[2 : semicolonPosition+2])
	if err != nil {
		return 0, fmt.Errorf("%s", err[0])
	}
	value := expression.Evaluate()
	environment.Global.Reset(variableName, value)
	return semicolonPosition + 3, nil
}

func handleReassignmentCondition(tokens []models.Token) (models.Node, error) {
	variableName := tokens[0].Lexem
	expression, err := parse.Parse(tokens[2:])
	if err != nil {
		return models.NilNode{}, fmt.Errorf("invalid reassignment expression")
	}
	value := expression.Evaluate()
	environment.Global.Reset(variableName, value)
	return expression, nil
}

func handleExpression(tokens []models.Token) (models.Node, error) {
	if utils.IsReassignmentCondition(tokens) {
		return handleReassignmentCondition(tokens)
	}
	if start, end, bool := utils.ExpressionHasFunctionCall(tokens); bool {
		result, _, err := handleFunCall(tokens[start : end+1])
		if err != nil {
			return models.NilNode{}, fmt.Errorf("invalid function call:%v", err.Error())
		}
		var value string
		switch v := result.Evaluate().(type) {
		case string:
			value = "\"" + v + "\""
		default:
			value = result.String()
		}
		funcTokens := tokenize.Tokenize(value)
		if len(funcTokens.Errors) != 0 {
			return models.NilNode{}, fmt.Errorf("error with tokenizing function call")
		}
		var newTokens []models.Token
		for i, token := range tokens {
			if i >= start && i <= end {
				if i == end {
					newTokens = append(newTokens, funcTokens.Success...)
				}
			} else {
				newTokens = append(newTokens, token)
			}
		}
		return handleExpression(newTokens)
	}
	expression, parseErrors := parse.Parse(tokens)
	if parseErrors != nil {
		return models.NilNode{}, fmt.Errorf("invalid expression: %v", parseErrors[0])
	}
	return expression, nil
}

func handlePrint(tokens []models.Token) (int, error) {
	if len(tokens) < 2 {
		return 0, fmt.Errorf("incomplete print statement")
	}
	semicolonPosition := utils.FindSemicolonPosition(tokens)
	if semicolonPosition == -1 {
		return 0, fmt.Errorf("no semicolon found after print")
	}
	expression, err := handleExpression(tokens[1:semicolonPosition])
	if err != nil {
		return 0, fmt.Errorf("invalid print expression: %v", err.Error())
	}
	result := expression.Evaluate()
	fmt.Printf("%v\n", result)
	return semicolonPosition + 1, nil
}

func handleIf(tokens []models.Token) (int, interface{}, error) {
	positions := findIfStatementPositions(tokens)
	if !positions.IsValid() {
		return 0, nil, fmt.Errorf("malformed if statement")
	}
	condition, err := handleExpression(tokens[positions.ConditionStart+1 : positions.ConditionEnd])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid if condition: %v", err.Error())
	}

	if condition.IsTruthy() {
		ret, err := Interprete(tokens[positions.IfBodyStart : positions.IfBodyEnd+1])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid if body: %v", err.Error())
		}
		if ret != nil {
			return positions.IfBodyEnd + 1, ret, err
		}
	} else {
		conditionMet := false
		for _, elseIfBlock := range positions.ElseIfBlocks {

			if elseIfBlock.BodyEnd > 0 && elseIfBlock.ConditionEnd > 0 && elseIfBlock.BodyStart > 0 && elseIfBlock.ConditionStart > 0 {
				elseIfCondition, err := handleExpression(tokens[elseIfBlock.ConditionStart+1 : elseIfBlock.ConditionEnd])
				if err != nil {
					return 0, nil, fmt.Errorf("invalid else-if condition: %v", err.Error())
				}
				if elseIfCondition.IsTruthy() {
					ret, err := Interprete(tokens[elseIfBlock.BodyStart : elseIfBlock.BodyEnd+1])

					if err != nil {
						return 0, nil, fmt.Errorf("invalid else-if body: %v", err.Error())
					}
					if ret != nil {
						return positions.IfBodyEnd + 1, ret, err
					}
					conditionMet = true
					break
				}
			}
		}

		if !conditionMet && positions.HasElseBlock() {
			ret, err := Interprete(tokens[positions.ElseBodyStart : positions.ElseBodyEnd+1])
			if err != nil {
				return 0, nil, fmt.Errorf("invalid else body: %v", err.Error())
			}
			if ret != nil {
				return positions.IfBodyEnd + 1, ret, err
			}
		}
	}

	if positions.HasElseBlock() {
		return positions.ElseBodyEnd + 1, nil, nil
	} else if len(positions.ElseIfBlocks) > 0 {
		return positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1].BodyEnd + 1, nil, nil
	}
	return positions.IfBodyEnd + 1, nil, nil
}

func findIfStatementPositions(tokens []models.Token) models.IfStatementPositions {
	positions := models.IfStatementPositions{
		ConditionStart: -1,
		ConditionEnd:   -1,
		IfBodyStart:    -1,
		IfBodyEnd:      -1,
		ElseIfBlocks:   make([]models.ElseIfBlock, 0),
		ElseBodyStart:  -1,
		ElseBodyEnd:    -1,
	}

	parenCount := 0
	braceCount := 0
	currentBlock := "if"
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Type {
		case models.LEFT_PAREN:
			if currentBlock == "if" && positions.ConditionStart == -1 && parenCount == 0 && braceCount == 0 {
				positions.ConditionStart = i
			} else if currentBlock == "elif" && parenCount == 0 && braceCount == 0 {
				positions.ElseIfBlocks = append(positions.ElseIfBlocks, models.ElseIfBlock{
					ConditionStart: i,
				})
			}
			parenCount++

		case models.RIGHT_PAREN:
			parenCount--
			if parenCount == 0 && braceCount == 0 {
				if currentBlock == "if" && positions.ConditionEnd == -1 {
					positions.ConditionEnd = i
					positions.IfBodyStart = i + 1
					if tokens[i+1].Type != models.LEFT_BRACE {
						positions.IfBodyEnd = utils.FindSemicolonPosition(tokens[i+1:]) + i + 1
					}
					if tokens[i+1].Type == models.IF {
						positions.IfBodyEnd = utils.FindLastSemicolonInSameLine(tokens[i+1:]) + i + 1
					}
				} else if currentBlock == "elif" && len(positions.ElseIfBlocks) > 0 {
					lastBlock := &positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1]
					lastBlock.ConditionEnd = i
					lastBlock.BodyStart = i + 1
					if tokens[i+1].Type != models.LEFT_BRACE {
						lastBlock.BodyEnd = utils.FindSemicolonPosition(tokens[i+1:]) + i + 1
					}
				}
			}

		case models.LEFT_BRACE:
			braceCount++
		case models.RIGHT_BRACE:
			braceCount--
			if braceCount == 0 && parenCount == 0 {
				if currentBlock == "if" && positions.IfBodyEnd == -1 {
					positions.IfBodyEnd = i
				} else if currentBlock == "elif" && len(positions.ElseIfBlocks) > 0 {
					positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1].BodyEnd = i
				} else {
					positions.ElseBodyEnd = i
					goto exit
				}
			}

		case models.SEMICOLON:
			if braceCount == 0 {
				if len(positions.ElseIfBlocks) > 0 && i == positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1].BodyEnd {
					if i+1 < len(tokens) && tokens[i+1].Type == models.ELSE {
						currentBlock = "else"
					} else {
						goto exit
					}
				}
				if i == positions.ElseBodyEnd {
					goto exit
				}
			}

		case models.ELSE:
			if braceCount == 0 {
				if i+1 < len(tokens) && tokens[i+1].Type == models.IF {
					currentBlock = "elif"
				} else {
					currentBlock = "else"
					if len(positions.ElseIfBlocks) > 0 {
						lastBlock := &positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1]
						if lastBlock.BodyEnd == -1 {
							lastBlock.BodyEnd = i - 1
						}
					} else if positions.IfBodyEnd == -1 {
						positions.IfBodyEnd = i - 1
					}

					if i+1 < len(tokens) && tokens[i+1].Type == models.LEFT_BRACE {
						positions.ElseBodyStart = i + 2
					} else {
						positions.ElseBodyStart = i + 1
						positions.ElseBodyEnd = utils.FindSemicolonPosition(tokens[i+1:]) + i + 1
					}
				}
			}

		}
	}
exit:
	return positions
}
