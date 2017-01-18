package gourd

type Context interface {
	Run(descr string, commands []string) (bool, error)
	And([]Expression) (bool, error)
	Or([]Expression) (bool, error)
}
