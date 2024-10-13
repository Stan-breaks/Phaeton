package parse

import (
	"strconv"
	"strings"

	"github.com/Stan-breaks/app/models"
)

func Parse(tokens models.Tokens) models.Node {
	left := 0
	right := 0
	if len(tokens.Success) == 3 {
		for index, item := range tokens.Success {
			splitToken := strings.Split(item, " ")
			if index == 0 && splitToken[0] == "NUMBER" {
				left, _ = strconv.Atoi(splitToken[1])
			} else if index == 2 && splitToken[0] == "NUMBER" {
				right, _ = strconv.Atoi(splitToken[1])
			}
		}
		return models.BinaryNode{
			Left:  models.NumberNode{Value: left},
			Op:    "+",
			Right: models.NumberNode{Value: right},
		}

	}
	return nil
}
