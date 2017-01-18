package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseProgram(t *testing.T) {

	tests := []struct {
		Name    string
		Program string
		Result  []UseCase
	}{
		{
			"Basic",
			`
Given a foo
	which foo

When I eat a foo
	foo eat

Then I get sick
	! system-status

Then I die
uh oh`,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			p := NewParser(strings.NewReader(test.Program))
			useCases, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.Result, useCases) {
				t.Errorf("expected %+v, got %+v", test.Result, useCases)
			}
		})
	}
}
