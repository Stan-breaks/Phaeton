package interpreter

import (
	"fmt"
	"strings"

	"github.com/Stan-breaks/app/environment"
	"github.com/Stan-breaks/app/models"
	"github.com/Stan-breaks/app/parse"
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
			tokensProcessed, err := handleReassignment(tokens[currentPosition:])
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
	semicolonPosition := findSemicolonPosition(tokens[3:])
	if semicolonPosition == -1 {
		return 0, fmt.Errorf("no semicolon found")
	}
	expression, err := parse.Parse(tokens[3 : semicolonPosition+3])
	if err != nil {
		return 0, fmt.Errorf("invalid assignment expression")
	}
	value := expression.Evaluate()
	environment.Environment[variableName] = value
	return semicolonPosition + 4, nil
}

func handleReassignment(tokens []models.TokenInfo) (int, error) {
	if !strings.HasPrefix(tokens[1].Token, "EQUAL") {
		return 0, fmt.Errorf("no equal found in reassignment")
	}
	variableName := strings.Split(tokens[0].Token, " ")[1]
	semicolonPosition := findSemicolonPosition(tokens[2:])
	if semicolonPosition == -1 {
		return 0, fmt.Errorf("no semicolon found in reassignment")
	}
	expression, err := parse.Parse(tokens[2 : semicolonPosition+2])
	if err != nil {
		return 0, fmt.Errorf("%s", err[0])
	}
	environment.Environment[variableName] = expression.Evaluate()
	return semicolonPosition + 3, nil
}

func handleReassignmentCondition(tokens []models.TokenInfo) (models.Node, error) {
	variableName := strings.Split(tokens[0].Token, " ")[1]
	expression, err := parse.Parse(tokens[2:])
	if err != nil {
		return models.NilNode{}, fmt.Errorf("invalid reassignment expression")
	}
	environment.Environment[variableName] = expression.Evaluate()
	return expression, nil
}

func handleExpression(tokens []models.TokenInfo) (models.Node, error) {
	if utils.IsReassignmentCondition(tokens) {
		return handleReassignmentCondition(tokens)
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
	semicolonPosition := findSemicolonPosition(tokens)
	if semicolonPosition == 0 {
		return 0, fmt.Errorf("no semicolon found after print")
	}
	expression, err := parse.Parse(tokens[1:semicolonPosition])
	if err != nil {
		return 0, fmt.Errorf("invalid print expression")
	}
	result := expression.Evaluate()
	fmt.Printf("%v\n", result)
	return semicolonPosition + 1, nil
}

func handleIf(tokens []models.TokenInfo) (int, error) {
	positions := findIfStatementPositions(tokens)
	if !positions.isValid() {
		return 0, fmt.Errorf("malformed if statement")
	}
	condition, err := handleExpression(tokens[positions.conditionStart+1 : positions.conditionEnd])
	if err != nil {
		return 0, fmt.Errorf("invalid if condition: %v", err.Error())
	}
	if condition.Evaluate().(bool) {
		err := Interprete(tokens[positions.ifBodyStart : positions.ifBodyEnd+1])
		if err != nil {
			return 0, fmt.Errorf("invalid if body: %v", err.Error())
		}
	} else if positions.hasElseBlock() {
		err := Interprete(tokens[positions.elseBodyStart : positions.elseBodyEnd+1])
		if err != nil {
			return 0, fmt.Errorf("invalid else body: %v", err.Error())
		}
	}

	if !positions.hasElseBlock() {
		return positions.ifBodyEnd + 1, nil
	}
	return positions.elseBodyEnd + 1, nil
}

type ifStatementPositions struct {
	conditionStart int
	conditionEnd   int
	ifBodyStart    int
	ifBodyEnd      int
	elseBodyStart  int
	elseBodyEnd    int
}

func (p ifStatementPositions) isValid() bool {
	return p.conditionStart != -1 && p.conditionEnd != -1 &&
		p.ifBodyStart != -1 && p.ifBodyEnd != -1
}

func (p ifStatementPositions) hasElseBlock() bool {
	return p.elseBodyEnd != -1 && p.elseBodyStart != -1
}

func findIfStatementPositions(tokens []models.TokenInfo) ifStatementPositions {
	positions := ifStatementPositions{
		conditionStart: -1,
		conditionEnd:   -1,
		ifBodyStart:    -1,
		ifBodyEnd:      -1,
		elseBodyStart:  -1,
		elseBodyEnd:    -1,
	}

	parenCount := 0
	braceCount := 0

	for i := 0; i < len(tokens); i++ {
		token := tokens[i].Token

		switch {
		case strings.HasPrefix(token, "LEFT_PAREN"):
			if positions.conditionStart == -1 && parenCount == 0 {
				positions.conditionStart = i
			}
			parenCount++

		case strings.HasPrefix(token, "RIGHT_PAREN"):
			parenCount--
			if parenCount == 0 && positions.conditionEnd == -1 {
				positions.conditionEnd = i
				positions.ifBodyStart = i + 1
				if !strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
					positions.ifBodyEnd = findSemicolonPosition(tokens[i+1:]) + i + 1
				}
			}

		case strings.HasPrefix(token, "LEFT_BRACE"):
			braceCount++

		case strings.HasPrefix(token, "RIGHT_BRACE"):
			braceCount--
			if braceCount == 0 {
				if positions.ifBodyEnd != -1 && positions.elseBodyEnd == -1 {
					positions.elseBodyEnd = i
				}
				if positions.ifBodyEnd == -1 {
					positions.ifBodyEnd = i
				}
			}

		case strings.HasPrefix(token, "ELSE"):
			if braceCount == 0 && positions.elseBodyStart == -1 {
				positions.ifBodyEnd = i - 1
				if strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
					positions.elseBodyStart = i + 2
				} else if strings.HasPrefix(tokens[i+1].Token, "IF") {
					continue
				} else {
					positions.elseBodyStart = i + 1
					positions.elseBodyEnd = findSemicolonPosition(tokens[i+1:]) + i + 1
				}
			}
		}
	}
	return positions
}

func findSemicolonPosition(tokens []models.TokenInfo) int {
	for i := 0; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i].Token, "SEMICOLON") {
			return i
		}
	}
	return -1
}
