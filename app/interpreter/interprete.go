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
	semicolonPosition := utils.FindSemicolonPosition(tokens[3:])
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
	semicolonPosition := utils.FindSemicolonPosition(tokens[2:])
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
	semicolonPosition := utils.FindSemicolonPosition(tokens)
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
	if !positions.IsValid() {
		return 0, fmt.Errorf("malformed if statement")
	}
	condition, err := handleExpression(tokens[positions.ConditionStart+1 : positions.ConditionEnd])
	if err != nil {
		return 0, fmt.Errorf("invalid if condition: %v", err.Error())
	}

	if condition.Evaluate().(bool) {
		err := Interprete(tokens[positions.IfBodyStart : positions.IfBodyEnd+1])
		if err != nil {
			return 0, fmt.Errorf("invalid if body: %v", err.Error())
		}
	} else {
		conditionMet := false
		for _, elseIfBlock := range positions.ElseIfBlocks {
			elseIfCondition, err := handleExpression(tokens[elseIfBlock.ConditionStart+1 : elseIfBlock.ConditionEnd])
			if err != nil {
				return 0, fmt.Errorf("invalid else-if condition: %v", err.Error())
			}
			if elseIfCondition.Evaluate().(bool) {
				err := Interprete(tokens[elseIfBlock.BodyStart : elseIfBlock.BodyEnd+1])
				if err != nil {
					return 0, fmt.Errorf("invalid else-if body: %v", err.Error())
				}
				conditionMet = true
				break
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
			if currentBlock == "if" && positions.ConditionStart == -1 && parenCount == 0 {
				positions.ConditionStart = i
			} else if currentBlock == "elif" && parenCount == 0 {
				positions.ElseIfBlocks = append(positions.ElseIfBlocks, models.ElseIfBlock{
					ConditionStart: i,
				})
			}
			parenCount++

		case strings.HasPrefix(token, "RIGHT_PAREN"):
			parenCount--
			if parenCount == 0 {
				if currentBlock == "if" && positions.ConditionEnd == -1 {
					positions.ConditionEnd = i
					positions.IfBodyStart = i + 1
					if !strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
						positions.IfBodyEnd = utils.FindSemicolonPosition(tokens[i+1:]) + i + 1
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
			if braceCount == 0 {
				if currentBlock == "if" && positions.IfBodyEnd == -1 {
					positions.IfBodyEnd = i
				} else if currentBlock == "elif" && len(positions.ElseIfBlocks) > 0 {
					positions.ElseIfBlocks[len(positions.ElseIfBlocks)-1].BodyEnd = i
					currentBlock = "if"
				} else if currentBlock == "else" {
					positions.ElseBodyEnd = i
				}
			}

		case strings.HasPrefix(token, "ELSE"):
			if braceCount == 0 {
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

					if strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
						positions.ElseBodyStart = i + 2
					} else {
						positions.ElseBodyStart = i + 1
						positions.ElseBodyEnd = utils.FindSemicolonPosition(tokens[i+1:]) + i + 1
					}
				}
			}
		}
	}

	return positions
}
