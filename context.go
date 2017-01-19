package gourd

// Context is an interface that defines how a parsed syntax tree is executed.
// Since Context is injected into the parsed tree to evaluate it, the behavior
// can easily be modified.
type Context interface {
	Run(descr string, commands []string) (bool, error)
	And([]Expression) (bool, error)
	Or([]Expression) (bool, error)
}
