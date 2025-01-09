package interpreter

import (
	"fmt"
	"github.com/Stan-breaks/app/models"
	"github.com/Stan-breaks/app/parse"
	"strings"
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
			ifLine := token.Line
			end := findMatchingEnd(ifLine, currentPosition, tokens)
			err := handleIf(tokens[currentPosition:end])
			if err != nil {
				return err
			}
			currentPosition = end

		case strings.HasPrefix(token.Token, "PRINT"):
			end, err := handlePrint(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += end
		case strings.HasPrefix(token.Token, "VAR"):

		default:
			currentPosition++
		}
	}
	return nil
}

func handlePrint(tokens []models.TokenInfo) (int, error) {
	if len(tokens) < 2 {
		return 0, fmt.Errorf("incomplete print statement")
	}
	expr, err := parse.Parse(tokens[1:])
	if err != nil {
		return 0, fmt.Errorf("invalid print expression")
	}
	result := expr.Evaluate()
	fmt.Print(result)
	tokensUsed := 2
	for i := 1; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i].Token, "SEMICOLON") {
			tokensUsed = i + 1
			break
		}
	}

	return tokensUsed, nil
}

func findMatchingEnd(initialLine int, currentPosition int, tokens []models.TokenInfo) int {
	braceCount := 0
	for i := currentPosition; i < len(tokens); i++ {
		token := tokens[i].Token
		if strings.HasPrefix(token, "LEFT_BRACE") {
			braceCount++
		} else if strings.HasPrefix(token, "RIGHT_BRACE") {
			braceCount--
			if braceCount == 0 && initialLine != tokens[i].Line {
				return i + 1
			}
		} else if strings.HasPrefix(token, "SEMICOLON") && braceCount == 0 {
			if initialLine == tokens[i].Line {
				return i + 1
			}
		}
	}
	return len(tokens)
}

func handleIf(tokens []models.TokenInfo) error {
	conditionStart := -1
	conditionEnd := -1
	bodyStart := -1
	bodyEnd := -1
	braceCount := 0

	for i := 0; i < len(tokens); i++ {
		token := tokens[i].Token

		switch {
		case strings.HasPrefix(token, "LEFT_PAREN") && conditionStart == -1:
			conditionStart = i
		case strings.HasPrefix(token, "RIGHT_PAREN"):
			conditionEnd = i
		case strings.HasPrefix(token, "LEFT_BRACE"):
			if bodyStart == -1 {
				bodyStart = i + 1
			}
			braceCount++
		case strings.HasPrefix(token, "RIGHT_BRACE"):
			braceCount--
			if braceCount == 0 {
				bodyEnd = i - 1
			}
		}
	}

	if bodyStart == -1 && conditionEnd != -1 {
		bodyStart = conditionEnd + 1
		bodyEnd = len(tokens) - 1
	}
	if conditionStart == -1 || conditionEnd == -1 || bodyStart == -1 {
		return fmt.Errorf("malformed if statement")
	}
	condition, err := parse.Parse(tokens[conditionStart : conditionEnd+1])
	if err != nil {
		return fmt.Errorf("invalid condition: %v", err)
	}

	if condition.Evaluate().(bool) {
		if bodyEnd == -1 {
			bodyEnd = len(tokens)
		}
		return Interprete(tokens[bodyStart:bodyEnd])
	}

	return nil
}
