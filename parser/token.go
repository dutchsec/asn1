package asn1parser

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS
	COMMENT

	// Literals
	IDENT // main

	// Misc characters
	OPTIONAL_TERM_OPEN  // [
	OPTIONAL_TERM_CLOSE // ]

	GROUP_OPEN  // {
	GROUP_CLOSE // }

	PARENTHESES_OPEN  // (
	PARENTHESES_CLOSE // )

	IMPORTS     // IMPORTS
	EXPORTS     // EXPORTS
	DEFINITIONS // DEFINITIONS
	FROM        // FROM
	BEGIN       // BEGIN
	END         // END
	APPLICATION // APPLICATION
	UNIVERSAL   // UNIVERSAL
	OPTIONAL    // OPTIONAL
	DEFAULT     // DEFAULT
	TRUE        // TRUE
	FALSE       // FALSE
	IMPLICIT    // IMPLICIT
	EXPLICIT    // EXPLICIT
	INTEGER     // INTEGER
	CHOICE      // CHOICE
	SET         // SET
	SEQUENCE    // SEQUENCE
	ENUMERATED  // ENUMERATED
	OCTET       // OCTET
	OF          // OF
	SIZE        // SIZE
	CHARACTER   // CHARACTER
	BIT         // BIT
	STRING      // STRING
	OBJECT      // OBJECT

	IDENTIFIER // IDENTIFIER

	COMMA      // ,
	SEMICOLON  // ;
	DOUBLE_DOT // ..
	TRIPLE_DOT // ...

	ASSIGNMENT_OPERATOR // ::=
	CONSTRAINED         // CONSTRAINED
	BY                  // BY

	VISIBLE_STRING    // VisibleString
	T61_STRING        // T61String
	PRINTABLE_STRING  // PrintableString
	OBJECT_DESCRIPTOR // ObjectDescriptor
	UTC_TIME          // UTCTime
	GENERALIZED_TIME  // GeneralizedTime
	NUMERIC_STRING    // NumericString
	UNIVERSAL_STRING  // UniversalString
	GENERAL_STRING    // GeneralString
	GRAPHIC_STRING    // GraphicString
)

// Token represents a lexical token.
type Token int

func (t Token) String() string {
	switch t {
	case ILLEGAL:
		return "<illegal>"
	case EOF:
		return "<eof>"
	case WS:
		return "<ws>"
	case COMMENT:
		return "<comment>"
	case IDENT:
		return "<ident>"
	case OPTIONAL_TERM_OPEN:
		return "<optional term open>"
	case OPTIONAL_TERM_CLOSE:
		return "<optional term close>"
	case GROUP_OPEN:
		return "<group open>"
	case GROUP_CLOSE:
		return "<group close>"
	case PARENTHESES_OPEN:
		return "<parentheses open>"
	case PARENTHESES_CLOSE:
		return "<parentheses close>"
	case SEMICOLON:
		return "<semicolon>"
	case COMMA:
		return "<comma>"
	case TRIPLE_DOT:
		return "<triple dot>"
	case ASSIGNMENT_OPERATOR:
		return "<assignment operator>"
	case EXPORTS:
		return "<exports>"
	case IMPORTS:
		return "<imports>"
	case BEGIN:
		return "<begin>"
	case END:
		return "<end>"
	case DEFINITIONS:
		return "<definitions>"
	case BY:
		return "<by>"
	case CONSTRAINED:
		return "<constrained>"
	}

	return "<invalid>"
}
