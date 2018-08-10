package asn1parser_test

import (
	"strings"
	"testing"

	"github.com/dutchsec/asn1-scheme-parser"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok asn1parser.Token
		lit string
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: asn1parser.EOF},
		{s: `#`, tok: asn1parser.ILLEGAL, lit: `#`},
		{s: ` `, tok: asn1parser.WS, lit: " "},
		{s: "\t", tok: asn1parser.WS, lit: "\t"},
		{s: "\n", tok: asn1parser.WS, lit: "\n"},

		// Misc characters
		{s: `;`, tok: asn1parser.SEMICOLON, lit: ";"},

		// Identifiers
		{s: `foo`, tok: asn1parser.IDENT, lit: `foo`},
		{s: `Zx12_3U_-`, tok: asn1parser.IDENT, lit: `Zx12_3U_-`},

		// Keywords
		{s: `DEFINITIONS`, tok: asn1parser.DEFINITIONS, lit: "DEFINITIONS"},
		{s: `IMPORTS`, tok: asn1parser.IMPORTS, lit: "IMPORTS"},
		{s: `EXPORTS`, tok: asn1parser.EXPORTS, lit: "EXPORTS"},
	}

	for i, tt := range tests {
		s := asn1parser.NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
