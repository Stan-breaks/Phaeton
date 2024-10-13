package models

import (
	"fmt"
	"strconv"

	"github.com/Stan-breaks/app/utils"
)

type Node interface {
	String() string
	Evaluate() interface{}
}

type NumberNode struct {
	Value int
}

func (n NumberNode) String() string {
	return strconv.Itoa(n.Value)
}

func (n NumberNode) Evaluate() interface{} {
	return n.Value
}

type BooleanNode struct {
	Value bool
}

func (n BooleanNode) String() string {
	return fmt.Sprint(n.Value)
}

func (n BooleanNode) Evaluate() interface{} {
	return n.Value
}

type BinaryNode struct {
	Left  Node
	Op    string
	Right Node
}

func (n BinaryNode) String() string {
	return fmt.Sprintf("%v %s %v", n.Left, n.Op, n.Right)
}

func (n BinaryNode) Evaluate() interface{} {
	left := n.Left.Evaluate()
	right := n.Right.Evaluate()
	switch n.Op {
	case string(utils.PLUS):
		return left.(int) + right.(int)
	case string(utils.MINUS):
		return left.(int) - right.(int)
	case string(utils.STAR):
		return left.(int) * right.(int)
	case string(utils.SLASH):
		return left.(int) / right.(int)
	case "and":
		return left.(bool) && right.(bool)
	case "or":
		return left.(bool) || right.(bool)
	default:
		panic("Unknown operator: " + n.Op)
	}

}
