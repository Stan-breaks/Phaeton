package interpreter

import (
	"fmt"
	"strings"

	"github.com/Stan-breaks/app/models"
	"github.com/Stan-breaks/app/parse"
)

func Interprete(tokens []models.TokenInfo) error {
	currentPosition := 0
	for currentPosition < len(tokens) {
		if strings.HasPrefix(tokens[currentPosition].Token, "IF") {
			ifLine := tokens[currentPosition].Line
			end := findMatchingEnd(ifLine, currentPosition, tokens)
			err := handleIf(tokens[currentPosition:end])
			if err != nil {
				return err
			}
			currentPosition = end
		} else if strings.HasPrefix(tokens[currentPosition].Token, "PRINT") {
			splitToken := strings.Split(tokens[currentPosition+1].Token, " ")
			fmt.Print(splitToken[2])
			currentPosition += 2
		} else {
			currentPosition++
		}
	}
	return nil
}

func findMatchingEnd(initialLine int, currentPosition int, tokens []models.TokenInfo) int {
	endToken := 0
	for i := currentPosition; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i].Token, "SEMICOLON") {
			if initialLine == tokens[i].Line {
				endToken = i
			} else {
				endToken = i + 1
			}
			return endToken
		}
	}
	return endToken
}

func handleIf(tokens []models.TokenInfo) error {
	currentPosition := 1
	startCondition := 0
	endCondition := 0
	startBody := 0
	endBody := 0
	for i := currentPosition; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i].Token, "LEFT_PAREN") {
			startCondition = i
		}
		if strings.HasPrefix(tokens[i].Token, "RIGHT_PAREN") {
			endCondition = i
			if strings.HasPrefix(tokens[i+1].Token, "LEFT_BRACE") {
				startBody = i + 2
			} else {
				startBody = i + 1
			}
		}
		if strings.HasPrefix(tokens[i].Token, "RIGHT_BRACE") {
			endBody = i
		}
	}
	if endBody == 0 {
		endBody = len(tokens) - 1
	}
	if startCondition != 0 && endCondition != 0 && startBody != 0 && endBody != 0 {
		return fmt.Errorf("invalid if")
	} else {
		condition, err := parse.Parse(tokens[startCondition:endCondition])
		if err != nil {
			return fmt.Errorf("invalid if condition")
		}
		if condition.Evaluate().(bool) {
			err := Interprete(tokens[startBody:endBody])
			if err != nil {
				return fmt.Errorf("invalid if body")
			}
		}
		return nil
	}
}
