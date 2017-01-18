package gourd_test

import (
	"strings"
	"testing"

	"github.com/jncornett/gourd"
)

var ExampleSpec = `
A bank account is a financial account maintained by a financial
institution for a customer.

Given: I have a bank account
	ls ./account.csv

And: The balance is 0

	! bank check-balance ./account.csv

When: I try to withdraw some money

	bank withdraw 1000 ./account.csv > out || true

Then: nothing will happen

	! grep '.*' out

(finally)

	rm -rf out
`

func TestLexer(t *testing.T) {
	tests := []struct {
		name   string
		spec   string
		tokens []gourd.TokenTuple
	}{
		{
			name: "Example",
			spec: ExampleSpec,
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex := gourd.NewLexer(strings.NewReader(test.spec))
			for i, tt := range test.tokens {
				var result *gourd.TokenTuple
				for {
					var err error
					result, err = lex.Lex()
					if err != nil {
						t.Fatalf("(#%v) %v", i, err)
						return
					}
					if result == nil {
						t.Fatalf("(#%v) Lex returned a nil token", i)
						return
					}
					// For brevity, filter out the text tokens
					if result.Token != gourd.TokenText {
						break
					}
				}
				if tt.Token != result.Token {
					t.Errorf("(#%v) Token mismatch: expected %v, got %v", i, tt, result)
				}
				if tt.Value != result.Value {
					t.Errorf("(#%v) Value mismatch: expected %v, got %v", i, tt, result)
				}
			}
		})
	}
}
