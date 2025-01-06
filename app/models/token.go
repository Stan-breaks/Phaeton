package models

type Tokens struct {
	Success []TokenInfo
	Errors  []string
}

type TokenInfo struct {
	Token string
	Line  int
}

type Token struct {
	Type  []string
	Value any
}
