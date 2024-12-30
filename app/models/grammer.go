package models

import (
	"fmt"

	"github.com/Stan-breaks/app/utils"
)

type Node interface {
	String() string
	Evaluate() interface{}
}

type NumberNode struct {
	Value float32
}

func (n NumberNode) String() string {
	return utils.FormatFloat(n.Value)
}

func (n NumberNode) Evaluate() interface{} {
	return n.Value
}

type StringNode struct {
	Value string
}

func (n StringNode) String() string {
	return n.Value
}

func (n StringNode) Evaluate() interface{} {
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

type NilNode struct {
}

func (n NilNode) String() string {
	return "nil"
}

func (n NilNode) Evaluate() interface{} {
	return "nil"
}

type ParenthesisNode struct {
	Expression Node
}

func (n ParenthesisNode) String() string {
	return fmt.Sprintf("(group %s)", n.Expression.String())
}

func (n ParenthesisNode) Evaluate() interface{} {
	return n.Expression.Evaluate()
}

type UnaryNode struct {
	Op    string
	Value Node
}

func (n UnaryNode) String() string {
	return fmt.Sprintf("(%s %v)", n.Op, n.Value)
}

func (n UnaryNode) Evaluate() interface{} {
	num := n.Value.Evaluate()
	switch n.Op {
	case string(utils.PLUS):
		return 1 * num.(float32)
	case string(utils.MINUS):
		return -1 * num.(float32)
	default:
		panic("Unknown operator: " + n.Op)
	}
}

type BinaryNode struct {
	Left  Node
	Op    string
	Right Node
}

func (n BinaryNode) String() string {
	return fmt.Sprintf("(%s %v %v)", n.Op, n.Left, n.Right)
}

func (n BinaryNode) Evaluate() interface{} {
	left := n.Left.Evaluate()
	right := n.Right.Evaluate()
	switch n.Op {
	case string(utils.PLUS):
		return left.(float32) + right.(float32)
	case string(utils.MINUS):
		return left.(float32) - right.(float32)
	case string(utils.STAR):
		return left.(float32) * right.(float32)
	case string(utils.SLASH):
		return left.(float32) / right.(float32)
	case "==":
		return left.(float32) == right.(float32)
	case "and":
		return left.(bool) && right.(bool)
	case "or":
		return left.(bool) || right.(bool)
	default:
		panic("Unknown operator: " + n.Op)
	}

}
