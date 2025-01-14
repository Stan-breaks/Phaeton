package utils

import (
	"strings"

	"github.com/Stan-breaks/app/models"
)

func IsReassignment(tokens []models.TokenInfo) bool {
	return strings.HasPrefix(tokens[0].Token, "IDENTIFIER") && strings.HasPrefix(tokens[1].Token, "EQUALS")
}
