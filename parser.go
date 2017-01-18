package gourd

type Parser interface {
	Parse() ([]Scenario, error)
}

type parser struct {
	lex Lexer
}

func NewParser(lex Lexer) Parser {
	return &parser{lex: lex}
}

func (p *parser) Parse() ([]Scenario, error) {
	return nil, nil
}
