package main

type UseCase struct {
	Givens    []Condition
	WhenThens []WhenThen
}

type Condition struct {
	Descr    string
	Commands []string
}

type WhenThen struct {
	When Condition
	Then Condition
}
