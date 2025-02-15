package utils

import (
	"strings"

	"github.com/Stan-breaks/app/models"
)

func FindNoOfArgs(tokens []models.TokenInfo) [][]models.TokenInfo {
	var result [][]models.TokenInfo
	var arr []models.TokenInfo
	var empty []models.TokenInfo
	if len(tokens) == 1 {
		arr = append(arr, tokens...)
		result = append(result, arr)
		return result
	}
	for _, token := range tokens {
		if strings.HasPrefix(token.Token, "COMMA") {
			result = append(result, arr)
			arr = empty
		} else {
			arr = append(arr, token)
		}
	}
	if len(arr) != 0 {
		result = append(result, arr)
	}
	return result
}

func IsFunctionCall(tokens []models.TokenInfo) bool {
	return strings.HasPrefix(tokens[1].Token, "LEFT_PAREN") && strings.HasPrefix(tokens[len(tokens)-2].Token, "RIGHT_PAREN")
}
func IsFunctionCallExpression(tokens []models.TokenInfo) bool {
	return strings.HasPrefix(tokens[1].Token, "LEFT_PAREN") && strings.HasPrefix(tokens[len(tokens)-1].Token, "RIGHT_PAREN")
}

func ExpressionHasFunctionCall(tokens []models.TokenInfo) (int, int, bool) {
	identifier := -1
	for i, token := range tokens {
		switch {
		case strings.HasPrefix(token.Token, "IDENTIFIER") && identifier == -1:
			identifier = i
		case strings.HasPrefix(token.Token, "RIGHT_PAREN") && identifier != -1:
			return identifier, i, IsFunctionCallExpression(tokens[identifier : i+1])
		}
	}
	return -1, -1, false
}
func IsReassignmentCondition(tokens []models.TokenInfo) bool {
	if len(tokens) < 2 {
		return false
	}
	return strings.HasPrefix(tokens[0].Token, "IDENTIFIER") && strings.HasPrefix(tokens[1].Token, "EQUAL")
}

func FindSemicolonPosition(tokens []models.TokenInfo) int {
	parenCount := 0
	braceCount := 0
	for i, token := range tokens {
		switch {
		case strings.HasPrefix(token.Token, "LEFT_PAREN"):
			parenCount++
		case strings.HasPrefix(token.Token, "RIGHT_PAREN"):
			parenCount--
		case strings.HasPrefix(token.Token, "LEFT_BRACE"):
			braceCount++
		case strings.HasPrefix(token.Token, "RIGHT_BRACE"):
			braceCount--
		case strings.HasPrefix(token.Token, "SEMICOLON"):
			if parenCount == 0 && braceCount == 0 {
				return i
			}
		}
	}
	return -1
}

func FindLastSemicolonInSameLine(tokens []models.TokenInfo) int {
	parenCount := 0
	braceCount := 0
	val := -1
	line := tokens[0].Line
	for i, token := range tokens {
		if token.Line != line {
			goto exit
		}
		switch {
		case strings.HasPrefix(token.Token, "LEFT_PAREN"):
			parenCount++
		case strings.HasPrefix(token.Token, "RIGHT_PAREN"):
			parenCount--
		case strings.HasPrefix(token.Token, "LEFT_BRACE"):
			braceCount++
		case strings.HasPrefix(token.Token, "RIGHT_BRACE"):
			braceCount--
		case strings.HasPrefix(token.Token, "SEMICOLON"):
			if parenCount == 0 && braceCount == 0 {
				val = i
			}
		}
	}
exit:
	return val
}

func FindClosingParen(tokens []models.TokenInfo) int {
	parenCount := 0
	for i, token := range tokens {
		switch {
		case strings.HasPrefix(token.Token, "LEFT_PAREN"):
			parenCount++
		case strings.HasPrefix(token.Token, "RIGHT_PAREN"):
			parenCount--
			if parenCount == 0 {
				return i
			}
		}
	}
	return 0
}
