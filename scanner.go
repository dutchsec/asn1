package asn1parser

import (
	"bufio"
	"bytes"
	"io"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isComment(ch) {
		// peek next char is comment also
		s.unread()
		return s.scanComment()
	} else if isDigit(ch) {
		s.unread()
		return s.scanIdent()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	if ch == '.' {
		for {
			if s.read() != '.' {
				s.unread()
				s.unread()
				break
			}

			if s.read() != '.' {
				s.unread()
				return DOUBLE_DOT, ""
			}

			return TRIPLE_DOT, ""
		}
	}

	if ch == ':' {
		for {
			if s.read() != ':' {
				s.unread()
				s.unread()
				break
			}

			if s.read() != '=' {
				s.unread()
				s.unread()
				s.unread()
				break
			}

			return ASSIGNMENT_OPERATOR, ""
		}
	}
	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '(':
		return PARENTHESES_OPEN, string(ch)
	case ')':
		return PARENTHESES_CLOSE, string(ch)
	case '[':
		return OPTIONAL_TERM_OPEN, string(ch)
	case ']':
		return OPTIONAL_TERM_CLOSE, string(ch)
	case '{':
		return GROUP_OPEN, string(ch)
	case '}':
		return GROUP_CLOSE, string(ch)
	case ';':
		return SEMICOLON, string(ch)
	case ',':
		return COMMA, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanComment consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanComment() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '\n' {
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return COMMENT, buf.String()
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
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' && ch != '-' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch buf.String() {
	case "CONSTRAINED":
		return CONSTRAINED, buf.String()
	case "BY":
		return BY, buf.String()
	case "BEGIN":
		return BEGIN, buf.String()
	case "END":
		return END, buf.String()
	case "FROM":
		return FROM, buf.String()
	case "DEFINITIONS":
		return DEFINITIONS, buf.String()
	case "IMPORTS":
		return IMPORTS, buf.String()
	case "EXPORTS":
		return EXPORTS, buf.String()
	case "UNIVERSAL":
		return UNIVERSAL, buf.String()
	case "ENUMERATED":
		return ENUMERATED, buf.String()
	case "APPLICATION":
		return APPLICATION, buf.String()
	case "OPTIONAL":
		return OPTIONAL, buf.String()
	case "DEFAULT":
		return DEFAULT, buf.String()
	case "TRUE":
		return TRUE, buf.String()
	case "FALSE":
		return FALSE, buf.String()
	case "IMPLICIT":
		return IMPLICIT, buf.String()
	case "EXPLICIT":
		return EXPLICIT, buf.String()
	case "OF":
		return OF, buf.String()
	case "INTEGER":
		return INTEGER, buf.String()
	case "SEQUENCE":
		return SEQUENCE, buf.String()
	case "STRING":
		return STRING, buf.String()
	case "BIT":
		return BIT, buf.String()
	case "IDENTIFIER":
		return IDENTIFIER, buf.String()
	case "UTCTime":
		return UTC_TIME, buf.String()
	case "ObjectDescriptor":
		return OBJECT_DESCRIPTOR, buf.String()
	case "GaphicString":
		return GRAPHIC_STRING, buf.String()
	case "VisibleString":
		return VISIBLE_STRING, buf.String()
	case "PrintableString":
		return PRINTABLE_STRING, buf.String()
	case "T61String":
		return T61_STRING, buf.String()
	case "OBJECT":
		return OBJECT, buf.String()
	case "OCTET":
		return OCTET, buf.String()
	case "SIZE":
		return SIZE, buf.String()
	case "SET":
		return SET, buf.String()
	case "CHOICE":
		return CHOICE, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isComment returns true if the rune is a comment.
func isComment(ch rune) bool { return (ch == '-') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
