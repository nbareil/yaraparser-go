package yaraparser

type Token int

const (
	// Special tokens
	Illegal Token = iota
	Eof
	WS

	Ident
	Comment
	// rule definition
	Rule
	Colon
	CurlyBraceOpen
	CurlyBraceClose

	// String definition
	VarIdentifier
	DoubleQuote
	QuotedString
	Equal
	QuestionMark

	// Condition definition
	ParenOpen
	ParenClose
	All
	Of
	Them
	Num
	Alpha
	Asterisk
)
