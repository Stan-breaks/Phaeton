package models

import (
	"fmt"

	"github.com/Stan-breaks/app/utils/format"
	"github.com/Stan-breaks/app/utils/runes"
)

type Node interface {
	String() string
	Evaluate() interface{}
}

type NumberNode struct {
	Value float32
}

func (n NumberNode) String() string {
	return format.FormatFloat(n.Value)
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
	case string(runes.PLUS):
		return 1 * num.(float32)
	case string(runes.MINUS):
		return -1 * num.(float32)
	case string(runes.BANG):
		switch num := num.(type) {
		case bool:
			return !num
		default:
			panic(fmt.Sprintf("Invalid type for ! operator: %T", num))
		}
	default:
		panic("Unknown operator: " + n.Op)
	}
}

type BinaryNode struct {
	Left    Node
	Op      string
	Right   Node
	Shifted int8
}

func (n BinaryNode) String() string {
	if n.Shifted != 0 {
		return fmt.Sprintf("(%s %v %v)", n.Op, n.Right, n.Left)
	} else {
		return fmt.Sprintf("(%s %v %v)", n.Op, n.Left, n.Right)
	}
}

func (n BinaryNode) Evaluate() interface{} {
	left := n.Left.Evaluate()
	right := n.Right.Evaluate()
	switch n.Op {
	case string(runes.PLUS):
		return left.(float32) + right.(float32)
	case string(runes.MINUS):
		return left.(float32) - right.(float32)
	case string(runes.STAR):
		return left.(float32) * right.(float32)
	case string(runes.SLASH):
		return left.(float32) / right.(float32)
	case "==":
		return left == right
	case "!=":
		return left != right
	case string(runes.GREATER):
		return left.(float32) > right.(float32)
	case string(runes.LESS):
		return left.(float32) < right.(float32)
	case ">=":
		return left.(float32) >= right.(float32)
	case "<=":
		return left.(float32) <= right.(float32)

	default:
		panic("Unknown operator: " + n.Op)
	}

}
