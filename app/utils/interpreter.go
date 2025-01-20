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
	for i := 0; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i].Token, "SEMICOLON") {
			return i
		}
	}
	return -1
}
