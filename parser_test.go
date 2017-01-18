package gourd_test

import (
	"reflect"
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
		return &lex.list[lex.i], nil
	}
	return nil, nil
}

func TestParser(t *testing.T) {
	tests := []struct {
		name      string
		tokens    []gourd.TokenTuple
		scenarios []gourd.Scenario
	}{
		{
			name: "Example",
			tokens: []gourd.TokenTuple{
				{Token: gourd.TokenGiven, Value: "I have a bank account"},
				{Token: gourd.TokenBlock, Value: "ls ./account.csv"},
				{Token: gourd.TokenAnd, Value: "The balance is 0"},
				{Token: gourd.TokenBlock, Value: "! bank check-balance ./account.csv"},
				{Token: gourd.TokenWhen, Value: "I try to withdraw some money"},
				{Token: gourd.TokenBlock, Value: "bank withdraw 1000 ./account.csv > out || true"},
				{Token: gourd.TokenThen, Value: "nothing will happen"},
				{Token: gourd.TokenBlock, Value: "! grep '.*' out"},
				{Token: gourd.TokenFinally},
				{Token: gourd.TokenBlock, Value: "rm -rf out"},
			},
			scenarios: []gourd.Scenario{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex := listLexer{list: test.tokens}
			p := gourd.NewParser(&lex)
			scenarios, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			for i, sc := range test.scenarios {
				if i >= len(scenarios) {
					t.Fatalf("(#%v) missing scenario", i)
				}
				if !reflect.DeepEqual(sc, scenarios[i]) {
					t.Errorf("(#%v) expected %v, got %v", i, sc, scenarios[i])
				}
			}
		})
	}
}
