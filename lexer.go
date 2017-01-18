package gourd

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jncornett/scan"
)

// Token is the set of lexical tokens for the gourd BDD language
type Token int

// The list of tokens.
const (
	// special
	TokenText Token = iota
	TokenBlock
	// keywords
	TokenGiven
	TokenAnd
	TokenOr
	TokenWhen
	TokenThen
	TokenFinally
)

var tokens = [...]string{
	TokenText:    "<Text>",
	TokenBlock:   "<Block>",
	TokenGiven:   "Given:",
	TokenAnd:     "And:",
	TokenOr:      "Or:",
	TokenWhen:    "When:",
	TokenThen:    "Then:",
	TokenFinally: "(finally)",
}

// String returns the string corresponding to token tok.
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// TokenTuple wraps a token and an associated string value together
type TokenTuple struct {
	Token
	Value string
}

func (tt TokenTuple) String() string {
	return fmt.Sprintf("(%v)%q", tt.Token, tt.Value)
}

// Lexer is an interface that wraps the Lex method.
// Lexers hold internal state and repeated calls to the Lex
// method return the next token read from the backing text.
type Lexer interface {
	Lex() (*TokenTuple, error)
}

type lexer struct {
	s          scan.BufferedScanner
	lastIndent int
}

// NewLexer returns a Lexer that reads from r.
// Repeated calls to the Lex method will tokenize the data from r.
func NewLexer(r io.Reader) Lexer {
	var s scan.Scanner
	s = bufio.NewScanner(r)
	s = newEmptyLineFilter(s)
	return &lexer{s: scan.NewBufferedScanner(s)}
}

func (lex *lexer) Lex() (*TokenTuple, error) {
	for lex.s.Scan() {
		line := lex.s.Text()
		indent := getIndent(line)
		if indent > 0 && indent >= lex.lastIndent {
			lex.lastIndent = indent
			return &TokenTuple{Token: TokenBlock, Value: trimPrefix(line, "")}, nil
		}
		tt := parseDirective(line)
		if tt != nil {
			return tt, nil
		}
		lex.lastIndent = 0
	}
	err := lex.s.Err()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

var directives = [...]Token{
	TokenGiven,
	TokenWhen,
	TokenThen,
	TokenAnd,
	TokenOr,
	TokenFinally,
}

func parseDirective(line string) *TokenTuple {
	for _, tok := range directives {
		prefix := tok.String()
		if strings.HasPrefix(line, prefix) {
			return &TokenTuple{Token: tok, Value: trimPrefix(line, prefix)}
		}
	}
	return nil
}

func newEmptyLineFilter(s scan.Scanner) scan.Scanner {
	return &scan.FilterScanner{
		Scanner: s,
		Filter: func(v scan.View) bool {
			for _, b := range v.Bytes() {
				if b != ' ' && b != '\t' {
					return true
				}
			}
			return false
		},
	}
}

func getIndent(line string) int {
	var (
		i  int
		ch rune
	)
	for i, ch = range line {
		if ch != ' ' && ch != '\t' {
			break
		}
	}
	return i
}

func trimPrefix(s, prefix string) string {
	return strings.Trim(strings.TrimPrefix(s, prefix), " \t")
}
