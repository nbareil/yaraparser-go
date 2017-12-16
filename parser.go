package yaraparser

import (
	"fmt"
	"io"
	"strings"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace/non-comment token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	if tok == Comment {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) parseSectionIdentifier() (string, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != Ident {
		return lit, fmt.Errorf("found %q, expecting a section name", lit)
	}
	section := strings.ToLower(lit)

	if tok, lit := p.scan(); tok != Colon {
		return "", fmt.Errorf("found %q, expecting ':'", lit)
	}

	return section, nil
}

// Parse parses
func (p *Parser) Parse() (*YaraRule, error) {
	rule := &YaraRule{}

	var tok Token
	lit := ""
	if tok, lit = p.scanIgnoreWhitespace(); tok != Rule {
		return nil, fmt.Errorf(`found %q, expected keyword "rule"`, lit)
	}

	if tok, lit = p.scanIgnoreWhitespace(); tok != Ident {
		return nil, fmt.Errorf("found %q, expected a valid rule identifier", lit)
	}
	rule.Name = lit

	if err := p.parseTags(rule); err != nil {
		return nil, err
	}
	if tok, lit = p.scanIgnoreWhitespace(); tok != CurlyBraceOpen {
		return nil, fmt.Errorf(`found %q after rulename, expecting "{"`, lit)
	}

	for {
		tok, lit = p.scanIgnoreWhitespace()
		p.unscan()

		if tok == CurlyBraceClose || tok == Eof {
			break
		}
		if err := p.parseSection(rule); err != nil {
			return nil, err
		}
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok == Eof {
		return nil, fmt.Errorf("found EOF, expecting '}'")
	}

	if tok != CurlyBraceClose {
		return nil, fmt.Errorf("found %q, expecing '}'", lit)
	}

	return rule, nil
}

func (p *Parser) parseSection(rule *YaraRule) error {
	section, err := p.parseSectionIdentifier()
	if err != nil {
		return err
	}

	switch section {
	case "meta":
		return p.parseMetaSection(rule)
	case "strings":
		return p.parseStringsSection(rule)
	case "conditions":
		return p.parseConditionsSection(rule)
	}

	return nil
}

func (p *Parser) parseMetaSection(rule *YaraRule) error {
	rule.Metas = map[string]string{}
	for {
		tok, key := p.scanIgnoreWhitespace()
		if tok == CurlyBraceClose || tok == Eof {
			p.unscan()
			break
		}

		if tok != Ident {
			return fmt.Errorf("found %q, expecting a key", key)
		}

		tok, lit := p.scanIgnoreWhitespace()
		if tok != Equal {
			return fmt.Errorf("found %q, expecting a value assignment", lit)
		}

		tok, val := p.scanIgnoreWhitespace()
		if tok != QuotedString {
			return fmt.Errorf("found %q, expecting a quoted string, token=%q", val, tok)
		}

		rule.Metas[key] = val
	}

	return nil
}

func (p *Parser) parseStringsSection(rule *YaraRule) error {
	rule.Strings = map[string]YaraPattern{}

	for {
		tok, name := p.scanIgnoreWhitespace()
		if tok == CurlyBraceClose || tok == Eof {
			p.unscan()
			break
		}

		if tok != VarIdentifier {
			return fmt.Errorf("found %q, expecting a variable identifier", name)
		}

		tok, lit := p.scanIgnoreWhitespace()
		if tok != Equal {
			return fmt.Errorf("found %q, expecting a value assignment", lit)
		}

		tok, val := p.scanIgnoreWhitespace()
		switch tok {
		case QuotedString:
			rule.Strings[name] = YaraPattern{
				Type:   RegularString,
				String: val,
			}
		}
		if tok != QuotedString {
			return fmt.Errorf("found %q, expecting a quoted string, token=%q", val, tok)
		}

	}
	return nil
}

func (p *Parser) parseConditionsSection(rule *YaraRule) error {
	return nil
}

func (p *Parser) parseTags(rule *YaraRule) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != Colon {
		p.unscan()
		// tags are optionnal
		return nil
	}
	for {
		tok, lit = p.scanIgnoreWhitespace()
		if tok == Eof || tok != Ident {
			break
		}
		rule.Tags = append(rule.Tags, lit)
	}

	if len(rule.Tags) == 0 {
		return fmt.Errorf("invalid tag name")
	}

	return nil
}
