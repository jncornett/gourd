package main

type Command interface{}

type commandList struct{}

func newCommandList(commands []string) Command {
	return &commandList{} // TODO implement
}
