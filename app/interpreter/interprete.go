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
			end, err := handleAssignment(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += end
		case strings.HasPrefix(token.Token, "IDENTIFIER"):
			end, err := handleReassignment(tokens[currentPosition:])
			if err != nil {
				return err
			}
			currentPosition += end
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
	valName := strings.Split(nameToken.Token, " ")[1]
	if !strings.HasPrefix(tokens[2].Token, "EQUAL") {
		return 0, fmt.Errorf("equal not found")
	}
	end := 0
	for i := 3; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i].Token, "SEMICOLON") {
			end = i
			break
		}
	}
	if end == 0 {
		return 0, fmt.Errorf("no semicolon found")
	}
	expr, err := parse.Parse(tokens[3:end])
	if err != nil {
		return 0, fmt.Errorf("invalid assignment expression")
	}
	value := expr.Evaluate()

	environment.Environment[valName] = value
	return end + 1, nil
}

func handleReassignment(tokens []models.TokenInfo) (int, error) {
	if !strings.HasPrefix(tokens[1].Token, "EQUAL") {
		fmt.Print(tokens)
		return 0, fmt.Errorf("no equal found in reassignment")
	}
	valname := strings.Split(tokens[0].Token, " ")[1]
	end := 0
	for i := 2; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i].Token, "SEMICOLON") {
			end = i
			break
		}
	}
	if end == 0 {
		return 0, fmt.Errorf("no semicolon found in reassignment")
	}
	val, err := parse.Parse(tokens[2:end])
	if err != nil {
		return 0, fmt.Errorf("%s", err[0])
	}
	environment.Environment[valname] = val.Evaluate()
	return 0, nil
}

func handleReassignmentCondition(tokens []models.TokenInfo) (models.Node, error) {
	valname := strings.Split(tokens[0].Token, " ")[1]
	val, err := parse.Parse(tokens[2:])
	if err != nil {
		return models.NilNode{}, fmt.Errorf("invalid reassignment expression")
	}
	environment.Environment[valname] = val.Evaluate()
	return val, nil
}

func handleExpression(tokens []models.TokenInfo) (models.Node, error) {
	var val models.Node
	var err error
	var error []string
	if utils.IsReassignmentCondition(tokens) {
		val, err = handleReassignmentCondition(tokens)
		if err != nil {
			return models.NilNode{}, err
		}
		return val, nil
	}
	val, error = parse.Parse(tokens)
	if error != nil {
		return models.NilNode{}, fmt.Errorf("invalid expression: %v", error[0])
	}
	return val, nil
}

func handlePrint(tokens []models.TokenInfo) (int, error) {
	if len(tokens) < 2 {
		return 0, fmt.Errorf("incomplete print statement")
	}

	tokensUsed := 0
	for i := 1; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i].Token, "SEMICOLON") {
			tokensUsed = i
			break
		}
	}
	expr, err := parse.Parse(tokens[1:tokensUsed])
	if err != nil {
		return 0, fmt.Errorf("invalid print expression")
	}
	result := expr.Evaluate()
	fmt.Print(result)

	if tokensUsed == 0 {
		return 0, fmt.Errorf("no semicolon found after print")
	}
	return tokensUsed + 2, nil
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
	firstRightParen := -1
	firstSemicolon := -1
	for i := 0; i < len(tokens); i++ {
		token := tokens[i].Token

		switch {
		case strings.HasPrefix(token, "LEFT_PAREN") && conditionStart == -1:
			conditionStart = i
		case strings.HasPrefix(token, "RIGHT_PAREN") && firstRightParen == -1:
			firstRightParen = i
		case strings.HasPrefix(token, "LEFT_BRACE"):
			conditionEnd = i - 1
			if bodyStart == -1 {
				bodyStart = i + 1
			}
			braceCount++
		case strings.HasPrefix(token, "RIGHT_BRACE"):
			braceCount--
			if braceCount == 0 {
				bodyEnd = i - 1
				goto exit
			}
		case strings.HasPrefix(token, "SEMICOLON") && firstSemicolon == -1:
			firstSemicolon = i
		}
	}
exit:
	if bodyStart == -1 && conditionEnd == -1 {
		bodyStart = firstRightParen + 1
		conditionEnd = firstRightParen
		bodyEnd = firstSemicolon
	}
	if conditionStart == -1 || conditionEnd == -1 || bodyStart == -1 {
		return fmt.Errorf("malformed if statement")
	}
	condition, err := handleExpression(tokens[conditionStart+1 : conditionEnd])
	if err != nil {
		return fmt.Errorf("invalid if expression: %v", err.Error())
	}

	if condition.Evaluate().(bool) {
		if bodyEnd == -1 {
			bodyEnd = len(tokens)
		}
		return Interprete(tokens[bodyStart : bodyEnd+1])
	}

	return nil
}
