package gourd

type Behavior struct {
	When Expression
	Then Expression
}

type Scenario struct {
	Given     Expression
	Behaviors []Behavior
}
