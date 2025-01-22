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
