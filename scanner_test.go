package yaraparser_test

import (
	"strings"
	"testing"

	"github.com/nbareil/yaraparser-go"
)

func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok yaraparser.Token
		lit string
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: yaraparser.Eof},
		{s: `#`, tok: yaraparser.Illegal, lit: `#`},
		{s: ` `, tok: yaraparser.WS, lit: " "},
		{s: "\t", tok: yaraparser.WS, lit: "\t"},
		{s: "\n", tok: yaraparser.WS, lit: "\n"},

		// single character
		{s: ":", tok: yaraparser.Colon, lit: ":"},
		{s: "=", tok: yaraparser.Equal, lit: "="},
		//{s: `"`, tok: yaraparser.DoubleQuote, lit: `"`},
		{s: "*", tok: yaraparser.Asterisk, lit: "*"},
		{s: "(", tok: yaraparser.ParenOpen, lit: "("},
		{s: ")", tok: yaraparser.ParenClose, lit: ")"},
		{s: "{", tok: yaraparser.CurlyBraceOpen, lit: "{"},
		{s: "}", tok: yaraparser.CurlyBraceClose, lit: "}"},

		// Modifier
		{s: "$0xblabla", tok: yaraparser.VarIdentifier, lit: "$0xblabla"},
		{s: `"foobar"`, tok: yaraparser.QuotedString, lit: "foobar"},
		{s: `" foo bar "`, tok: yaraparser.QuotedString, lit: " foo bar "},

		// keywords
		{s: "rule", tok: yaraparser.Rule, lit: "rule"},
		{s: "all", tok: yaraparser.All, lit: "all"},
		{s: "of", tok: yaraparser.Of, lit: "of"},
		{s: "them", tok: yaraparser.Them, lit: "them"},

		// misc
		{s: "foobar", tok: yaraparser.Ident, lit: "foobar"},
	}

	for i, tt := range tests {
		s := yaraparser.NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
