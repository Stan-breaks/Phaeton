package models

import (
	"fmt"
	"strconv"
)

type Node interface {
	String() string
}

type NumberNode struct {
	Value int
}

func (n NumberNode) String() string {
	return strconv.Itoa(n.Value)
}

type VariableNode struct {
	Name string
}

type BooleanNode struct {
	Value bool
}

func (n BooleanNode) String() string {
	return fmt.Sprint(n.Value)
}

type BinaryNode struct {
	Left  Node
	Op    string
	Right Node
}

func (n BinaryNode) String() string {
	return fmt.Sprintf("%v %s %v", n.Left, n.Op, n.Right)
}
