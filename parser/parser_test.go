package asn1parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dutchsec/asn1/parser"
)

func Example() {
	r := strings.NewReader(`
DEFINITIONS ::=

BEGIN
END
`)

	parser := asn1parser.NewParser(r)

	def, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	_ = def
}

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
	var tests = []struct {
		s   string
		def *asn1parser.ASNDefinition
		err string
	}{
		{s: `
-- ASN definition from
-- http://www.sisconet.com/techinfo.htm
-- slightly modified
--
--
--Corrections made July 2, 1994
--
--
-- Modified to pass asn2wrs

MMS { iso standard 9506 part(2) mms-general-module-version(2) }

DEFINITIONS ::=

BEGIN
END
`, def: &asn1parser.ASNDefinition{
			Name:  "MMS",
			Types: []asn1parser.ASNType{},
		},
			err: ""},
	}

	for i, tt := range tests {
		def, err := asn1parser.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.def, def) {
			t.Errorf("%d. %q\n\ndefinition mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.def, def)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
