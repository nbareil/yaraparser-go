package yaraparser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isDoubleQuote(ch rune) bool {
	return ch == '"'
}

func isVarIdentifier(ch rune) bool {
	return ch == '$'
}

func isTagMarker(ch rune) bool {
	return ch == ':'
}

func isHexFormat(ch rune) bool {
	return isDigit(ch) || ch == '?' || ch == '*'
}

// Scanner is bla
type Scanner struct {
	r *bufio.Reader
}

// NewScanner creates a new Scanner instance
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if ch == '"' {
		return s.scanQuotedString()
	} else if ch == '/' {
		return s.scanComment()
	} else if isVarIdentifier(ch) {
		s.unread()
		return s.scanVarIdentifier()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return Eof, ""
	case ':':
		return Colon, ":"
	case '{':
		return CurlyBraceOpen, "{"
	case '}':
		return CurlyBraceClose, "}"
	case '=':
		return Equal, "="
	case '*':
		return Asterisk, "*"
	case '(':
		return ParenOpen, "("
	case ')':
		return ParenClose, ")"
	}

	return Illegal, string(ch)
}

func (s *Scanner) scanVarIdentifier() (tok Token, lit string) {
	var buf bytes.Buffer

	buf.WriteRune(s.read())

	for {
		ch := s.read()

		if ch == eof {
			break
		} else if !isDigit(ch) && !isLetter(ch) || ch == '_' {
			s.unread()
			break
		}
		buf.WriteRune(ch)
	}

	return VarIdentifier, buf.String()
}
func (s *Scanner) scanQuotedString() (tok Token, lit string) {
	var buf bytes.Buffer
	var prev rune

	for {
		ch := s.read()
		if ch == eof {
			break
		} else if ch == '"' && prev != '\\' {
			// closing quote, do not capture it
			break
		}

		buf.WriteRune(ch)
		prev = ch
	}
	return QuotedString, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	kwd := strings.ToLower(buf.String())
	switch kwd {
	case "rule":
		return Rule, kwd
	case "all":
		return All, kwd
	case "of":
		return Of, kwd
	case "them":
		return Them, kwd
	}

	// Otherwise return as a regular identifier.
	return Ident, buf.String()
}

func (s *Scanner) scanComment() (tok Token, lit string) {
	var buf bytes.Buffer

	ch := s.read()
	if ch == eof {
		s.unread()
		return Illegal, ""
	}
	if ch == '/' {
		for {
			ch := s.read()
			if ch == eof {
				break
			} else if ch == '\n' {
				break
			} else if ch == '\r' {
				ch := s.read()
				if ch == '\n' || ch == eof {
					break
				}
			}
			buf.WriteRune(ch)
		}
	} else if ch == '*' {
		for {
			ch := s.read()
			if ch == eof {
				break
			}
			if ch == '/' {
				break
			}
			buf.WriteRune(ch)
		}
		buf.UnreadRune() // removes *
	}
	return Comment, buf.String()
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break

		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}
