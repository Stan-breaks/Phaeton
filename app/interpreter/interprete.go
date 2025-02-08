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

func Interprete(tokens []models.TokenInfo) error {
	currentPosition := 0
	for currentPosition < len(tokens) {
		if currentPosition >= len(tokens) {
			break
		}
		token := tokens[currentPosition]
		switch {
		case strings.HasPrefix(token.Token, "IF"):
			tokensProcessed, err := handleIf(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += tokensProcessed
		case strings.HasPrefix(token.Token, "PRINT"):
			tokensProcessed, err := handlePrint(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += tokensProcessed
		case strings.HasPrefix(token.Token, "VAR"):
			tokensProcessed, err := handleAssignment(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += tokensProcessed
		case strings.HasPrefix(token.Token, "IDENTIFIER"):
			if utils.IsFunctionCall(tokens[currentPosition:]) {
				tokensProcessed, err := handleFunCall(tokens[currentPosition:])
				if err != nil {
					return err
				}
				currentPosition += tokensProcessed
			} else {
				tokensProcessed, err := handleReassignment(tokens[currentPosition:])
				if err != nil {
					return err
				}
				currentPosition += tokensProcessed
			}

		case strings.HasPrefix(token.Token, "LEFT_PAREN"):
			tokensProcessed, err := handleParenStatement(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += tokensProcessed
		case strings.HasPrefix(token.Token, "WHILE"):
			tokensProcessed, err := handleWhile(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += tokensProcessed
		case strings.HasPrefix(token.Token, "FOR"):
			tokensProcessed, err := handleFor(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += tokensProcessed
		case strings.HasPrefix(token.Token, "FUN"):
			tokensProcessed, err := handleFun(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += tokensProcessed
		case strings.HasPrefix(token.Token, "RETURN"):
			tokensProcessed, err := handleReturn(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += tokensProcessed
		default:
			currentPosition++
		}
	}
	return nil
}

func handleReturn(tokens []models.TokenInfo) (int, error) {
	semicolon := utils.FindLastSemicolonInSameLine(tokens)
	if semicolon == -1 {
		return 0, fmt.Errorf("no semicolon in return statement")
	}
	result, err := parse.Parse(tokens[1:])
	if len(err) != 0 {
		return 0, fmt.Errorf("error with parsing return statement")
	}
	environment.Global.SetReturn(result.Evaluate())
	return semicolon + 1, nil
}

func handleExprFunCall(tokens []models.TokenInfo) (models.Node, int, error) {
	environment.Global.PushScope()
	defer environment.Global.PopScope()
	funName := strings.Split(tokens[0].Token, " ")[1]
	value, bool := environment.Global.Get(funName)
	if !bool {
		return models.NilNode{}, 0, fmt.Errorf("function not defined")
	}
	switch v := value.(type) {
	case models.Function:
		switch a := v.Arguments.(type) {
		case []models.TokenInfo:
			argumentEnd := len(tokens) - 1
			args := tokens[2:argumentEnd]
			if len(a) != len(args) {
				return models.NilNode{}, 0, fmt.Errorf("invalid no of function argument")
			}
			if len(a) > 0 {
				for i := 0; i < len(a); i++ {
					if !strings.HasPrefix(a[i].Token, "COMMA") {
						valName := strings.Split(a[i].Token, " ")[1]
						value, err := parse.Parse(args[i : i+1])

						if err != nil {
							return models.NilNode{}, 0, fmt.Errorf("invalid arguments")
						}
						environment.Global.Set(valName, value.Evaluate())
					}
				}

			}
		}
		switch b := v.Body.(type) {
		case []models.TokenInfo:
			err := Interprete(b)
			if err != nil {
				return models.NilNode{}, 0, err
			}
		}

	}
	return models.NilNode{}, len(tokens), nil
}

func handleFunCall(tokens []models.TokenInfo) (int, error) {
	environment.Global.PushScope()
	defer environment.Global.PopScope()
	semicolon := utils.FindSemicolonPosition(tokens)
	if semicolon == -1 {
		return 0, fmt.Errorf("no semicolon found")
	}
	funName := strings.Split(tokens[0].Token, " ")[1]
	value, bool := environment.Global.Get(funName)
	if !bool {
		return 0, fmt.Errorf("function not defined")
	}
	switch v := value.(type) {
	case models.Function:
		switch a := v.Arguments.(type) {
		case []models.TokenInfo:
			argumentEnd := utils.FindClosingParen(tokens)
			args := tokens[2:argumentEnd]
			if len(a) != len(args) {
				return 0, fmt.Errorf("invalid no of function argument")
			}
			if len(a) > 0 {
				for i := 0; i < len(a); i++ {
					if !strings.HasPrefix(a[i].Token, "COMMA") {
						valName := strings.Split(a[i].Token, " ")[1]
						value, err := parse.Parse(args[i : i+1])

						if err != nil {
							return 0, fmt.Errorf("invalid arguments")
						}
						environment.Global.Set(valName, value.Evaluate())
					}
				}

			}
		}
		switch b := v.Body.(type) {
		case []models.TokenInfo:
			err := Interprete(b)
			if err != nil {
				return 0, err
			}
		}

	}
	return semicolon + 1, nil
}

func handleFun(tokens []models.TokenInfo) (int, error) {
	positions := findFunPositions(tokens)
	if !positions.IsValid() {
		return positions.BodyEnd + 1, fmt.Errorf("invalid function")
	}
	funName := strings.Split(tokens[1].Token, " ")[1]
	function := models.Function{
		Arguments: tokens[positions.ArgumentStart+1 : positions.ArgumentEnd],
		Body:      tokens[positions.BodyStart : positions.BodyEnd+1],
	}
	environment.Global.Set(funName, function)
	return positions.BodyEnd + 1, nil
}

func findFunPositions(tokens []models.TokenInfo) models.FunStatementPositions {
	positions := models.FunStatementPositions{
		ArgumentStart: -1,
		ArgumentEnd:   -1,
		BodyStart:     -1,
		BodyEnd:       -1,
	}
	parenCount := 0
	braceCount := 0
	for i := 0; i < len(tokens); i++ {
		switch {
		case strings.HasPrefix(tokens[i].Token, "LEFT_PAREN"):
			if parenCount == 0 && positions.ArgumentStart == -1 {
				positions.ArgumentStart = i
			}
			parenCount++
		case strings.HasPrefix(tokens[i].Token, "RIGHT_PAREN"):
			parenCount--
			if parenCount == 0 && positions.ArgumentEnd == -1 {
				positions.ArgumentEnd = i
				positions.BodyStart = i + 1

			}
		case strings.HasPrefix(tokens[i].Token, "LEFT_BRACE"):
			braceCount++
		case strings.HasPrefix(tokens[i].Token, "RIGHT_BRACE"):
			braceCount--
			if braceCount == 0 && positions.BodyEnd == -1 {
				positions.BodyEnd = i
			}
		}
	}
	return positions
}

func handleFor(tokens []models.TokenInfo) (int, error) {
	environment.Global.PushScope()
	defer environment.Global.PopScope()

	positions := findForPositions(tokens)
	if !positions.IsValid() {
		return positions.BodyEnd + 1, fmt.Errorf("invalid for statement")
	}
	declarationStart := positions.ConditionStart + 1
	declarationEnd := utils.FindSemicolonPosition(tokens[positions.ConditionStart+1:positions.ConditionEnd]) + positions.ConditionStart + 1
	conditionStart := declarationEnd + 1
	conditionEnd := utils.FindSemicolonPosition(tokens[declarationEnd+1:positions.ConditionEnd]) + declarationEnd + 1
	expressionStart := conditionEnd + 1
	expressionEnd := positions.ConditionEnd

	err := Interprete(tokens[declarationStart : declarationEnd+1])
	if err != nil {
		return positions.BodyEnd + 1, err
	}
	condition, err := handleExpression(tokens[conditionStart:conditionEnd])
	if err != nil {
		return positions.BodyEnd + 1, err
	}
	if condition.IsTruthy() {
		if expressionStart != expressionEnd {
			for {
				environment.Global.PushScope()
				err := Interprete(tokens[positions.BodyStart : positions.BodyEnd+1])
				environment.Global.PopScope()
				if err != nil {
					return positions.BodyEnd + 1, err
				}
				_, err = handleExpression(tokens[expressionStart:expressionEnd])
				if err != nil {
					return positions.BodyEnd + 1, err
				}

				condition, err = handleExpression(tokens[conditionStart:conditionEnd])
				if err != nil {
					return positions.BodyEnd + 1, err
				}
				if !condition.IsTruthy() {
					break
				}
			}
		} else {
			for {
				err := Interprete(tokens[positions.BodyStart : positions.BodyEnd+1])
				if err != nil {
					return positions.BodyEnd + 1, err
				}
				condition, err = handleExpression(tokens[conditionStart:conditionEnd])
				if err != nil {
					return positions.BodyEnd + 1, err
				}
				if !condition.IsTruthy() {
					break
				}
			}
		}
	}
	return positions.BodyEnd + 1, nil
}

func findForPositions(tokens []models.TokenInfo) models.ForStatementPositions {
	positions := models.ForStatementPositions{
		ConditionStart: 1,
		ConditionEnd:   -1,
		BodyStart:      -1,
		BodyEnd:        -1,
	}
	parenCount := 0
	braceCount := 0
	for i := 1; i < len(tokens); i++ {
		token := tokens[i].Token
		switch {
		case strings.HasPrefix(token, "LEFT_PAREN"):
			parenCount++
		case strings.HasPrefix(token, "RIGHT_PAREN"):
			parenCount--
			if parenCount == 0 && positions.ConditionEnd == -1 {
				positions.ConditionEnd = i
				positions.BodyStart = i + 1
				if !strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
					positions.BodyEnd = utils.FindLastSemicolonInSameLine(tokens[i+1:]) + i + 1
				}
			}
		case strings.HasPrefix(token, "LEFT_BRACE"):
			braceCount++
		case strings.HasPrefix(token, "RIGHT_BRACE"):
			braceCount--
			if braceCount == 0 && positions.BodyEnd == -1 {
				positions.BodyEnd = i
				goto exit
			}
		case strings.HasPrefix(token, "SEMICOLON") && positions.BodyEnd == i && parenCount == 0 && braceCount == 0:
			goto exit
		}

	}
exit:
	return positions
}

func handleWhile(tokens []models.TokenInfo) (int, error) {
	environment.Global.PushScope()
	defer environment.Global.PopScope()
	positions := findWhilePositions(tokens)
	if !positions.IsValid() {
		return positions.BodyEnd + 1, fmt.Errorf("invalid while statement")
	}
	condition, err := handleExpression(tokens[positions.ConditionStart+1 : positions.ConditionEnd])
	if err != nil {
		return positions.BodyEnd + 1, fmt.Errorf("invalid while condition")
	}
	if condition.IsTruthy() {
		for {
			err := Interprete(tokens[positions.BodyStart : positions.BodyEnd+1])
			if err != nil {
				return positions.BodyEnd + 1, fmt.Errorf("invalid while body")
			}
			condition, err = handleExpression(tokens[positions.ConditionStart+1 : positions.ConditionEnd])
			if err != nil {
				return positions.BodyEnd + 1, fmt.Errorf("invalid while condition")
			}
			if !condition.IsTruthy() {
				break
			}
		}

	}
	return positions.BodyEnd + 1, nil
}

func findWhilePositions(tokens []models.TokenInfo) models.WhileStatementPositions {
	positions := models.WhileStatementPositions{
		ConditionStart: 1,
		ConditionEnd:   -1,
		BodyStart:      -1,
		BodyEnd:        -1,
	}
	parenCount := 0
	braceCount := 0
	for i := 1; i < len(tokens); i++ {
		token := tokens[i].Token
		switch {
		case strings.HasPrefix(token, "LEFT_PAREN"):
			parenCount++
		case strings.HasPrefix(token, "RIGHT_PAREN"):
			parenCount--
			if parenCount == 0 && positions.ConditionEnd == -1 {
				positions.ConditionEnd = i
				positions.BodyStart = i + 1
				if !strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
					positions.BodyEnd = utils.FindLastSemicolonInSameLine(tokens[i+1:]) + i + 1
				}
			}
		case strings.HasPrefix(token, "LEFT_BRACE"):
			braceCount++
		case strings.HasPrefix(token, "RIGHT_BRACE"):
			braceCount--
			if braceCount == 0 && positions.BodyEnd == -1 {
				positions.BodyEnd = i
				goto exit
			}
		case strings.HasPrefix(token, "SEMICOLON") && positions.BodyEnd == i && parenCount == 0 && braceCount == 0:
			goto exit
		}

	}
exit:
	return positions
}

func handleParenStatement(tokens []models.TokenInfo) (int, error) {
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
		if truthy && strings.HasPrefix(tokens[op].Token, "OR") {
			break
		} else if !truthy && strings.HasPrefix(tokens[op].Token, "AND") {
			break
		}
		startParen = endParen + 2
		endParen = utils.FindClosingParen(tokens[startParen:]) + startParen
	}
	return end, nil
}

func handleCondition(tokens []models.TokenInfo) (bool, error) {
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

func handleAssignment(tokens []models.TokenInfo) (int, error) {
	if len(tokens) < 4 {
		return 0, fmt.Errorf("incomplete variable declaration")
	}
	nameToken := tokens[1]
	if !strings.HasPrefix(nameToken.Token, "IDENTIFIER") {
		return 0, fmt.Errorf("no identifier found")
	}
	variableName := strings.Split(nameToken.Token, " ")[1]
	if !strings.HasPrefix(tokens[2].Token, "EQUAL") {
		return 0, fmt.Errorf("equal not found")
	}
	semicolonPosition := utils.FindSemicolonPosition(tokens[3:])
	if semicolonPosition == -1 {
		return 0, fmt.Errorf("no semicolon found")
	}
	expression, err := parse.Parse(tokens[3 : semicolonPosition+3])
	if err != nil {
		return 0, fmt.Errorf("invalid assignment expression")
	}
	value := expression.Evaluate()
	environment.Global.Set(variableName, value)
	return semicolonPosition + 4, nil
}

func handleReassignment(tokens []models.TokenInfo) (int, error) {
	if !strings.HasPrefix(tokens[1].Token, "EQUAL") {
		return 0, fmt.Errorf("no equal found in reassignment")
	}
	variableName := strings.Split(tokens[0].Token, " ")[1]
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

func handleReassignmentCondition(tokens []models.TokenInfo) (models.Node, error) {
	variableName := strings.Split(tokens[0].Token, " ")[1]
	expression, err := parse.Parse(tokens[2:])
	if err != nil {
		return models.NilNode{}, fmt.Errorf("invalid reassignment expression")
	}
	value := expression.Evaluate()
	environment.Global.Reset(variableName, value)
	return expression, nil
}

func handleExpression(tokens []models.TokenInfo) (models.Node, error) {
	if utils.IsReassignmentCondition(tokens) {
		return handleReassignmentCondition(tokens)
	}
	if start, end, bool := utils.ExpressionHasFunctionCall(tokens); bool {
		result, _, err := handleExprFunCall(tokens[start : end+1])
		fmt.Println(err)
		if err != nil {
			return models.NilNode{}, fmt.Errorf("invalid function call")
		}
		value := result.String()
		funcTokens := tokenize.Tokenize(value, len(value))
		if len(funcTokens.Errors) == 0 {
			return models.NilNode{}, fmt.Errorf("error with tokenizing function call")
		}
		var newTokens []models.TokenInfo
		for i, token := range tokens {
			if i >= start && i <= end {
				if i == end {
					newTokens = append(newTokens, funcTokens.Success...)
				}
			} else {
				newTokens = append(newTokens, token)
			}
		}
		expression, parseErrors := parse.Parse(newTokens)
		if parseErrors != nil {
			return models.NilNode{}, fmt.Errorf("invalid expression: %v", parseErrors[0])
		}
		return expression, nil
	}
	expression, parseErrors := parse.Parse(tokens)
	if parseErrors != nil {
		return models.NilNode{}, fmt.Errorf("invalid expression: %v", parseErrors[0])
	}
	return expression, nil
}

func handlePrint(tokens []models.TokenInfo) (int, error) {
	if len(tokens) < 2 {
		return 0, fmt.Errorf("incomplete print statement")
	}
	semicolonPosition := utils.FindSemicolonPosition(tokens)
	if semicolonPosition == -1 {
		return 0, fmt.Errorf("no semicolon found after print")
	}
	expression, err := handleExpression(tokens[1:semicolonPosition])
	if err != nil {
		return 0, fmt.Errorf("invalid print expression")
	}
	result := expression.Evaluate()
	fmt.Printf("%v\n", result)
	return semicolonPosition + 1, nil
}

func handleIf(tokens []models.TokenInfo) (int, error) {
	positions := findIfStatementPositions(tokens)
	if !positions.IsValid() {
		return 0, fmt.Errorf("malformed if statement")
	}
	condition, err := handleExpression(tokens[positions.ConditionStart+1 : positions.ConditionEnd])
	if err != nil {
		return 0, fmt.Errorf("invalid if condition: %v", err.Error())
	}

	if condition.IsTruthy() {
		err := Interprete(tokens[positions.IfBodyStart : positions.IfBodyEnd+1])
		if err != nil {
			return 0, fmt.Errorf("invalid if body: %v", err.Error())
		}
	} else {
		conditionMet := false
		for _, elseIfBlock := range positions.ElseIfBlocks {

			if elseIfBlock.BodyEnd > 0 && elseIfBlock.ConditionEnd > 0 && elseIfBlock.BodyStart > 0 && elseIfBlock.ConditionStart > 0 {
				elseIfCondition, err := handleExpression(tokens[elseIfBlock.ConditionStart+1 : elseIfBlock.ConditionEnd])
				if err != nil {
					return 0, fmt.Errorf("invalid else-if condition: %v", err.Error())
				}
				if elseIfCondition.IsTruthy() {
					err := Interprete(tokens[elseIfBlock.BodyStart : elseIfBlock.BodyEnd+1])

					if err != nil {
						return 0, fmt.Errorf("invalid else-if body: %v", err.Error())
					}
					conditionMet = true
					break
				}
			}
		}

		if !conditionMet && positions.HasElseBlock() {
			err := Interprete(tokens[positions.ElseBodyStart : positions.ElseBodyEnd+1])
			if err != nil {
				return 0, fmt.Errorf("invalid else body: %v", err.Error())
			}
		}
	}

	if positions.HasElseBlock() {
		return positions.ElseBodyEnd + 1, nil
	} else if len(positions.ElseIfBlocks) > 0 {
		return positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1].BodyEnd + 1, nil
	}
	return positions.IfBodyEnd + 1, nil
}

func findIfStatementPositions(tokens []models.TokenInfo) models.IfStatementPositions {
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
		token := tokens[i].Token
		switch {
		case strings.HasPrefix(token, "LEFT_PAREN"):
			if currentBlock == "if" && positions.ConditionStart == -1 && parenCount == 0 && braceCount == 0 {
				positions.ConditionStart = i
			} else if currentBlock == "elif" && parenCount == 0 && braceCount == 0 {
				positions.ElseIfBlocks = append(positions.ElseIfBlocks, models.ElseIfBlock{
					ConditionStart: i,
				})
			}
			parenCount++

		case strings.HasPrefix(token, "RIGHT_PAREN"):
			parenCount--
			if parenCount == 0 && braceCount == 0 {
				if currentBlock == "if" && positions.ConditionEnd == -1 {
					positions.ConditionEnd = i
					positions.IfBodyStart = i + 1
					if !strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
						positions.IfBodyEnd = utils.FindSemicolonPosition(tokens[i+1:]) + i + 1
					}
					if strings.HasPrefix(tokens[i+1].Token, "IF") {
						positions.IfBodyEnd = utils.FindLastSemicolonInSameLine(tokens[i+1:]) + i + 1
					}
				} else if currentBlock == "elif" && len(positions.ElseIfBlocks) > 0 {
					lastBlock := &positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1]
					lastBlock.ConditionEnd = i
					lastBlock.BodyStart = i + 1
					if !strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
						lastBlock.BodyEnd = utils.FindSemicolonPosition(tokens[i+1:]) + i + 1
					}
				}
			}

		case strings.HasPrefix(token, "LEFT_BRACE"):
			braceCount++
		case strings.HasPrefix(token, "RIGHT_BRACE"):
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

		case strings.HasPrefix(token, "SEMICOLON") && braceCount == 0:
			if len(positions.ElseIfBlocks) > 0 && i == positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1].BodyEnd {
				if i+1 < len(tokens) && strings.HasPrefix(tokens[i+1].Token, "ELSE") {
					currentBlock = "else"
				} else {
					goto exit
				}
			}
			if i == positions.ElseBodyEnd {
				goto exit
			}

		case strings.HasPrefix(token, "ELSE") && braceCount == 0:
			if i+1 < len(tokens) && strings.HasPrefix(tokens[i+1].Token, "IF") {
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

				if i+1 < len(tokens) && strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
					positions.ElseBodyStart = i + 2
				} else {
					positions.ElseBodyStart = i + 1
					positions.ElseBodyEnd = utils.FindSemicolonPosition(tokens[i+1:]) + i + 1
				}
			}
		}
	}
exit:
	return positions
}
