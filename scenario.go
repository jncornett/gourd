package gourd

import "fmt"

// Behavior models a behavior that a feature is supposed to have.
type Behavior struct {
	When Expression
	Then Expression
}

func (b Behavior) String() string {
	return fmt.Sprintf("Behavior{When: %v Then: %v}", b.When, b.Then)
}

// Scenario models a single testing scenario with a Given (preconditions),
// and zero or more Behaviors.
type Scenario struct {
	Given     Expression
	Behaviors []Behavior
}

func (s Scenario) String() string {
	return fmt.Sprintf("Scenario{%v %v}", s.Given, s.Behaviors)
}
