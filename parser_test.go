package gourd_test

import (
	"fmt"
	"testing"

	"github.com/jncornett/gourd"
)

type listLexer struct {
	list []gourd.TokenTuple
	i    int
}

func (lex *listLexer) Lex() (*gourd.TokenTuple, error) {
	if 0 <= lex.i && lex.i < len(lex.list) {
		lex.i++
		return &lex.list[lex.i-1], nil
	}
	return nil, nil
}

func TestParser(t *testing.T) {
	tests := []struct {
		name   string
		tokens []gourd.TokenTuple
		out    string
		err    bool
	}{
		// FIXME these test cases need to be added!
		{
			name: "Example",
			tokens: []gourd.TokenTuple{
				{Token: gourd.TokenGiven, Value: "I have a bank account"},
				{Token: gourd.TokenBlock, Value: "ls ./account.csv"},
				{Token: gourd.TokenAnd, Value: "The balance is 0"},
				{Token: gourd.TokenBlock, Value: "! bank check-balance ./account.csv"},
				{Token: gourd.TokenWhen, Value: "I try to withdraw some money"},
				{
					Token: gourd.TokenBlock,
					Value: "bank withdraw 1000 ./account.csv > out || true",
				},
				{Token: gourd.TokenThen, Value: "nothing will happen"},
				{Token: gourd.TokenBlock, Value: "! grep '.*' out"},
				{Token: gourd.TokenFinally},
				{Token: gourd.TokenBlock, Value: "rm -rf out"},
			},
			out: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex := listLexer{list: test.tokens}
			p := gourd.NewParser(&lex)
			scenarios, err := p.Parse()
			if test.err {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(scenarios) < 1 {
				t.Fatal("nil scenario")
			}
			scenario := scenarios[0]
			out := fmt.Sprint(scenario)
			if test.out != out {
				t.Errorf("expected out to be %q, got %q", test.out, out)
			}
		})
	}
}
