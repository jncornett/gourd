package gourd

// Parser is an interface that wraps the Parse method.
type Parser interface {
	Parse() ([]Scenario, error)
}

type parser struct {
	lex      Lexer
	last     *TokenTuple
	buffered bool
}

// NewParser returns a new Parser based on lex.
func NewParser(lex Lexer) Parser {
	return &parser{lex: lex}
}

func (p *parser) Parse() (out []Scenario, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for {
		if sc := p.parseScenario(); sc != nil {
			out = append(out, *sc)
		} else {
			break
		}
	}
	return
}

func (p *parser) parseScenario() *Scenario {
	given := p.parseChunk(TokenGiven)
	if given == nil {
		return nil
	}
	var behaviors []Behavior
	for {
		b := p.parseBehavior()
		if b == nil {
			break
		}
		behaviors = append(behaviors, *b)
	}
	return &Scenario{Given: given, Behaviors: behaviors}
}

func (p *parser) parseBehavior() *Behavior {
	when := p.parseChunk(TokenWhen)
	then := p.parseChunk(TokenThen)
	if when == nil && then == nil {
		return nil
	}
	return &Behavior{
		When: when,
		Then: then,
	}
}

func (p *parser) parseChunk(t Token) Expression {
	tt := p.scanIgnoreText()
	if tt == nil {
		return nil
	}
	if tt.Token != t {
		p.unscan()
		return nil
	}
	head := p.parseBlock(*tt)
	tt = p.scanIgnoreText()
	if tt != nil {
		if tt.Token == TokenAnd {
			return &AndExpression{Children: p.parseTrailingChunks(head, *tt)}
		} else if tt.Token == TokenOr {
			return &OrExpression{Children: p.parseTrailingChunks(head, *tt)}
		} else {
			p.unscan()
		}
	}
	return head
}

func (p *parser) parseTrailingChunks(
	head Expression,
	next TokenTuple,
) []Expression {
	out := []Expression{head, p.parseBlock(next)}
	for {
		tt := p.scanIgnoreText()
		if tt == nil {
			break
		}
		if tt.Token != next.Token {
			p.unscan()
			break
		}
		out = append(out, p.parseBlock(*tt))
	}
	return out
}

func (p *parser) parseBlock(head TokenTuple) Expression {
	var commands []string
	for {
		tt := p.scan()
		if tt == nil {
			break
		}
		if tt.Token != TokenBlock {
			p.unscan()
			break
		}
	}
	return &BlockExpression{
		Description: head.Value,
		Commands:    commands,
	}
}

func (p *parser) scanIgnoreText() (tt *TokenTuple) {
	for {
		if tt = p.scan(); tt == nil || tt.Token != TokenText {
			break
		}
	}
	return
}

func (p *parser) scan() *TokenTuple {
	if p.buffered {
		p.buffered = false
		return p.last
	}
	var err error
	p.last, err = p.lex.Lex()
	if err != nil {
		panic(err)
	}
	return p.last
}

func (p *parser) unscan() {
	p.buffered = true
}
