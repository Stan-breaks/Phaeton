package interpreter

import (
	"fmt"
	"strings"

	"github.com/Stan-breaks/app/models"
)

func Interprete(tokens []models.TokenInfo) {
	currentPosition := 0
	for currentPosition < len(tokens) {
		if strings.HasPrefix(tokens[currentPosition].Token, "IF") {
			ifLine := tokens[currentPosition].Line
			end := findMatchingEnd(ifLine, currentPosition, tokens)
			handleIf(tokens[currentPosition:end])
			currentPosition = end
		} else {
			currentPosition++
		}
	}
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

}
