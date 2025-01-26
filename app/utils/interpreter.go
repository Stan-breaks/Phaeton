package utils

import (
	"strings"

	"github.com/Stan-breaks/app/models"
)

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
