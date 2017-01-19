package gourd

import "fmt"

// Expression is an interface that wraps the Eval method
type Expression interface {
	Eval(Context) (bool, error)
}

// AndExpression represents an expression that should be
// evaluated by 'and'-ing together all of it's child expressions.
type AndExpression struct {
	Children []Expression
}

func (e AndExpression) String() string {
	return fmt.Sprintf("And{%v}", e.Children)
}

// Eval evaluates the expression based on ctx and returns
// the evaluation result.
func (e *AndExpression) Eval(ctx Context) (bool, error) {
	return ctx.And(e.Children)
}

// OrExpression represents an expression that should be
// evaluated by 'or'-ing together all of it's child expressions.
type OrExpression struct {
	Children []Expression
}

func (e OrExpression) String() string {
	return fmt.Sprintf("Or{%v}", e.Children)
}

// Eval evaluates the expression based on ctx and returns
// the evaluation result.
func (e *OrExpression) Eval(ctx Context) (bool, error) {
	return ctx.Or(e.Children)
}

// BlockExpression represents the most basic expression.
// It is a description coupled with a list of commands.
type BlockExpression struct {
	Description string
	Commands    []string
}

// Eval evaluates the expression based on ctx and returns
// the evaluation result.
func (e *BlockExpression) Eval(ctx Context) (bool, error) {
	return ctx.Run(e.Description, e.Commands)
}

func (e BlockExpression) String() string {
	return fmt.Sprintf("Block{%v %v}", e.Description, e.Commands)
}
