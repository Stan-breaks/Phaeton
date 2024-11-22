func parseMultipleBinaryExpr(tokens []string) models.Node {
	// Precedence map for operators
	precedence := map[string]int{
		"*":           2,
		"/":           2,
		"PLUS":        1,
		"MINUS":       1,
		"LESS":        0,
		"EQUAL_EQUAL": 0,
		"AND":         -1,
		"OR":          -2,
	}

	// Recursive descent parser with precedence
	var parse func([]string, int) models.Node
	parse = func(tokens []string, minPrecedence int) models.Node {
		// Base case: single value or parenthesized expression
		if len(tokens) == 1 {
			splitToken := strings.Split(tokens[0], " ")
			return parsevalue(splitToken)
		}

		// Find the lowest precedence operator from right to left
		var leftNode models.Node
		var remainingTokens []string

		for i := len(tokens) - 1; i >= 0; i-- {
			splitToken := strings.Split(tokens[i], " ")
			if len(splitToken) > 1 && isOperator(splitToken[0]) {
				op := parseOperator(splitToken)

				// Check if this operator's precedence is high enough
				if precedence[op] >= minPrecedence {
					// Split tokens
					leftTokens := tokens[:i]
					rightTokens := tokens[i+1:]

					// Recursively parse left and right sides
					leftNode = parse(leftTokens, precedence[op]+1)
					remainingTokens = tokens[:i]

					// Create binary node
					return models.BinaryNode{
						Left:  leftNode,
						Op:    op,
						Right: parse(rightTokens, precedence[op]),
					}
				}
			}
		}

		// If no operator found, parse as a single value or expression
		if len(tokens) == 1 {
			splitToken := strings.Split(tokens[0], " ")
			return parsevalue(splitToken)
		}

		return models.NilNode{}
	}

	return parse(tokens, -2) // Start with lowest precedence
}
