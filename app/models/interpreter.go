package models

type ElseIfBlock struct {
	ConditionStart int
	ConditionEnd   int
	BodyStart      int
	BodyEnd        int
}

type IfStatementPositions struct {
	ConditionStart int
	ConditionEnd   int
	IfBodyStart    int
	IfBodyEnd      int
	ElseIfBlocks   []ElseIfBlock
	ElseBodyStart  int
	ElseBodyEnd    int
}

func (p IfStatementPositions) IsValid() bool {
	return p.ConditionStart != -1 && p.ConditionEnd != -1 &&
		p.IfBodyStart != -1 && p.IfBodyEnd != -1
}

func (p IfStatementPositions) HasElseBlock() bool {
	return p.ElseBodyEnd != -1 && p.ElseBodyStart != -1
}
