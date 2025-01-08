package format

import (
	"fmt"
	"strings"
)

func FormatFloat(v float32) string {
	s := fmt.Sprintf("%.5f", v)
	s = strings.Trim(s, "0")
	if string(s[len(s)-1]) == "." {
		s += "0"
	}
	if string(s[0]) == "." {
		s = "0" + s
	}
	return s
}
