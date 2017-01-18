package gourd

type Expression interface {
	Eval(Context) (bool, error)
}

type AndExpression struct {
	Children []Expression
}

func (e *AndExpression) Eval(ctx Context) (bool, error) {
	return ctx.And(e.Children)
}

type OrExpression struct {
	Children []Expression
}

func (e *OrExpression) Eval(ctx Context) (bool, error) {
	return ctx.Or(e.Children)
}

type BlockExpression struct {
	Description string
	Commands    []string
}

func (e *BlockExpression) Eval(ctx Context) (bool, error) {
	return ctx.Run(e.Description, e.Commands)
}
