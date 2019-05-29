package ccmath

import (
	"errors"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/stack"
	"github.com/golang-collections/go-datastructures/queue"
)

type bxtNode struct {
	Left  *bxtNode
	Value string
	Right *bxtNode
}

// Solve given postfix expression
func Solve(pf string) (float64, error) {
	root, err := buildTree(pf)
	if err != nil {
		return 0, err
	}
	return root.solve(), nil
}

// buildTree builds a tree with the given postfix expression
func buildTree(pf string) (*bxtNode, error) {
	tf := strings.Split(pf, " ")
	q := queue.New(int64(len(tf)))
	for _, curr := range tf {
		if curr == "" {
			continue
		}
		q.Put(&bxtNode{
			Value: curr,
		})
	}
	s := stack.New()
	for !q.Empty() {
		nextt, err := q.Get(1)
		if err != nil {
			return nil, errors.New("Improper postfix expression")
		}
		next := nextt[0].(*bxtNode)
		if strings.Contains("+-/*", next.Value) {
			next.Right = s.Pop().(*bxtNode)
			next.Left = s.Pop().(*bxtNode)
			s.Push(next)
		} else {
			s.Push(next)
		}
	}
	return s.Pop().(*bxtNode), nil
}

func (b bxtNode) solve() float64 {
	if b.Value == "" {
		return 0
	}
	if !strings.Contains("+-/*", b.Value) {
		value, err := strconv.ParseFloat(b.Value, 64)
		if err != nil {
			return 0
		}
		return value
	}
	switch b.Value {
	case "+":
		return b.Left.solve() + b.Right.solve()
	case "-":
		return b.Left.solve() - b.Right.solve()
	case "/":
		return b.Left.solve() / b.Right.solve()
	case "*":
		return b.Left.solve() * b.Right.solve()
	}
	return 0
}
