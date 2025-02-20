package utils

import (
	"github.com/Stan-breaks/app/models"
)

func FindNoOfArgs(tokens []models.Token) [][]models.Token {
	var result [][]models.Token
	var arr []models.Token
	var empty []models.Token
	if len(tokens) == 1 {
		arr = append(arr, tokens...)
		result = append(result, arr)
		return result
	}
	for _, token := range tokens {
		if token.Type == models.COMMA {
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

func IsFunctionCall(tokens []models.Token) bool {
	return tokens[1].Type == models.LEFT_PAREN && tokens[len(tokens)-2].Type == models.RIGHT_PAREN
}
func IsFunctionCallExpression(tokens []models.Token) bool {
	return tokens[1].Type == models.LEFT_PAREN && tokens[len(tokens)-1].Type == models.RIGHT_PAREN
}

func ExpressionHasFunctionCall(tokens []models.Token) (int, int, bool) {
	identifier := -1
	for i, token := range tokens {
		switch {
		case token.Type == models.IDENTIFIER && identifier == -1:
			identifier = i
		case token.Type == models.RIGHT_PAREN && identifier != -1:
			return identifier, i, IsFunctionCallExpression(tokens[identifier : i+1])
		}
	}
	return -1, -1, false
}
func IsReassignmentCondition(tokens []models.Token) bool {
	if len(tokens) < 2 {
		return false
	}
	return tokens[0].Type == models.IDENTIFIER && tokens[1].Type == models.EQUAL
}

func FindSemicolonPosition(tokens []models.Token) int {
	parenCount := 0
	braceCount := 0
	for i, token := range tokens {
		switch token.Type {
		case models.LEFT_PAREN:
			parenCount++
		case models.RIGHT_PAREN:
			parenCount--
		case models.LEFT_BRACE:
			braceCount++
		case models.RIGHT_BRACE:
			braceCount--
		case models.SEMICOLON:
			if parenCount == 0 && braceCount == 0 {
				return i
			}
		}
	}
	return -1
}

func FindLastSemicolonInSameLine(tokens []models.Token) int {
	parenCount := 0
	braceCount := 0
	val := -1
	line := tokens[0].Line
	for i, token := range tokens {
		if token.Line != line {
			goto exit
		}
		switch token.Type {
		case models.LEFT_PAREN:
			parenCount++
		case models.RIGHT_PAREN:
			parenCount--
		case models.LEFT_BRACE:
			braceCount++
		case models.RIGHT_BRACE:
			braceCount--
		case models.SEMICOLON:
			if parenCount == 0 && braceCount == 0 {
				val = i
			}
		}
	}
exit:
	return val
}

func FindClosingParen(tokens []models.Token) int {
	parenCount := 0
	for i, token := range tokens {
		switch token.Type {
		case models.LEFT_PAREN:
			parenCount++
		case models.RIGHT_PAREN:
			parenCount--
			if parenCount == 0 {
				return i
			}
		}
	}
	return 0
}
