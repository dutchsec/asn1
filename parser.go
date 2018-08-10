package asn1parser

import (
	"fmt"
	"io"
)

// Parser represents a parser.
type Parser struct {
	s *Scanner

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

// Parse parses an ASN1 Definition.
func (p *Parser) Parse() (*ASNDefinition, error) {
	d := &ASNDefinition{
		Types: []ASNType{},
	}

	// name of definition
	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("parser: found %q, expected IDENT", lit)
	} else {
		d.Name = lit
	}

	// between brackets
	if tok, _ := p.scanIgnoreWhitespace(); tok == GROUP_OPEN {
		for {
			if tok, _ = p.scanIgnoreWhitespace(); tok == GROUP_CLOSE {
				break
			}
		}
	} else {
		p.unscan()
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != DEFINITIONS {
		return nil, fmt.Errorf("parser: found %q, expected DEFINITIONS identifier", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != ASSIGNMENT_OPERATOR {
		return nil, fmt.Errorf("parser: found %q, expected ASSIGNMENT_OPERATOR", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != BEGIN {
		return nil, fmt.Errorf("parser: found %q, expected BEGIN", lit)
	}

	// loop through all types
	for {
		if tok, _ := p.scanIgnoreWhitespace(); tok == EXPORTS {
			if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
				return nil, fmt.Errorf("exports: found %+v, expected IDENT: %+v", tok, lit)
			}

			for {
				tok, lit := p.scanIgnoreWhitespace()
				if tok == COMMA {
				} else if tok == SEMICOLON {
					break
				} else {
					_ = lit
				}
			}

			// TODO: implement EXPORTS
			continue
		} else {
			p.unscan()
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok == IMPORTS {
			for {
				if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
					return nil, fmt.Errorf("imports: found %+v, expected IDENT: %+v", tok, lit)
				} else {
					// TODO: implement IMPORTS
					_ = lit
				}

				if tok, _ = p.scanIgnoreWhitespace(); tok != COMMA {
					p.unscan()
					break
				}

			}

			if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
				return nil, fmt.Errorf("imports: found %+v, expected FROM: %+v", tok, lit)
			}

			if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
				return nil, fmt.Errorf("imports: found %+v, expected IDENT: %+v", tok, lit)
			}

			if tok, _ := p.scanIgnoreWhitespace(); tok != GROUP_OPEN {
				p.unscan()
			} else {
				for {
					if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
						p.unscan()
						return nil, fmt.Errorf("imports: found %q, expected IDENT", lit)
					} else {
					}

					if tok, lit := p.scanIgnoreWhitespace(); tok == PARENTHESES_OPEN {
						if tok, lit = p.scanIgnoreWhitespace(); tok != IDENT {
							p.unscan()
						} else {
						}

						if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
							return nil, fmt.Errorf("imports: found %q, expected close parentheses", lit)
						}
					} else {
						p.unscan()
					}

					if tok, _ := p.scanIgnoreWhitespace(); tok == GROUP_CLOSE {
						break
					} else {
						p.unscan()
					}
				}
			}

			if tok, lit := p.scanIgnoreWhitespace(); tok != SEMICOLON {
				return nil, fmt.Errorf("imports: found %+v, expected SEMICOLON: %+v", tok, lit)
			}
			continue
		} else {
			p.unscan()
		}

		cmmn := ASNCommon{}

		if tok, lit := p.scanIgnoreWhitespace(); tok == END {
			p.unscan()
			break
		} else if tok == IDENT {
			// NAME
			cmmn.SetName(lit)
		} else {
			return nil, fmt.Errorf("decl: found %+v, expected IDENT: %+v", tok, lit)
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok != OBJECT {
			p.unscan()
		} else if tok, _ := p.scanIgnoreWhitespace(); tok != IDENTIFIER {
			p.unscan()
			p.unscan()
		} else {
			if tok, lit := p.scanIgnoreWhitespace(); tok != ASSIGNMENT_OPERATOR {
				return nil, fmt.Errorf("decl: found %q, expected ASSIGNMENT_OPERATOR", lit)
			}

			if tok, lit := p.scanIgnoreWhitespace(); tok != GROUP_OPEN {
				p.unscan()
				return nil, fmt.Errorf("found %q, expected GROUP_OPEN", lit)
			}

			for {
				if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
					p.unscan()
					return nil, fmt.Errorf("found %q, expected IDENT", lit)
				} else {
				}

				if tok, lit := p.scanIgnoreWhitespace(); tok == PARENTHESES_OPEN {
					if tok, lit = p.scanIgnoreWhitespace(); tok != IDENT {
						p.unscan()
					} else {
					}

					if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
						return nil, fmt.Errorf("found %q, expected close parentheses", lit)
					}
				} else {
					p.unscan()
				}

				if tok, _ := p.scanIgnoreWhitespace(); tok == GROUP_CLOSE {
					break
				} else {
					p.unscan()
				}
			}

			continue
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok == IDENT {
			alias := ASNAlias{
				ASNCommon: cmmn,
				Alias:     lit,
			}

			if tok, lit := p.scanIgnoreWhitespace(); tok != ASSIGNMENT_OPERATOR {
				return nil, fmt.Errorf("found %q, expected ASSIGNMENT_OPERATOR", lit)
			}

			if tok, lit := p.scanIgnoreWhitespace(); tok == IDENT {
				alias.Default = lit
			} else {
				p.unscan()
			}

			d.Types = append(d.Types, &alias)
		} else {
			p.unscan()

			if tok, lit := p.scanIgnoreWhitespace(); tok != ASSIGNMENT_OPERATOR {
				return nil, fmt.Errorf("found %q, expected ASSIGNMENT_OPERATOR", lit)
			}

			if tok, _ := p.scanIgnoreWhitespace(); tok != OPTIONAL_TERM_OPEN {
				p.unscan()
			} else {
				// [UNIVERSAL 8]
				if tok, lit := p.scanIgnoreWhitespace(); tok == APPLICATION {
				} else if tok == UNIVERSAL {
				} else {
					return nil, fmt.Errorf("found %q, expected IDENT", lit)
				}

				if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
					return nil, fmt.Errorf("found %q, expected IDENT", lit)
				}

				if tok, lit := p.scanIgnoreWhitespace(); tok != OPTIONAL_TERM_CLOSE {
					return nil, fmt.Errorf("found %q, expected OPTIONAL_TERM_CLOSE", lit)
				}
			}

			if tok, lit = p.scanIgnoreWhitespace(); tok == IMPLICIT {
				// IMPLICIT
				// TODO
				cmmn.Implicit = true
			} else {
				p.unscan()
			}

			if type_, err := p.scanType(cmmn); err != nil {
				return nil, err
			} else {
				d.Types = append(d.Types, type_)
			}

		}
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != END {
		return nil, fmt.Errorf("found %q, expected END", lit)
	}

	// Return the successfully parsed definition.
	return d, nil
}

func (p *Parser) scanSequence(cmmn ASNCommon) (ASNType, error) {
	// TODO: should we differentiate between ASNSequence and ASNSequenceOf?
	sequence := &ASNSequence{
		cmmn,
		"",
		nil,
	}

	if tok, _ := p.scanIgnoreWhitespace(); tok != SIZE {
		p.unscan()
	} else if tok, lit := p.scanIgnoreWhitespace(); tok == PARENTHESES_OPEN {
		if _, _, err := p.scanRange(); err != nil {
			return nil, fmt.Errorf("scanSequence: found %q, expected range", lit)
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			return nil, fmt.Errorf("scanSequence: found %q, expected comma", lit)
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok != TRIPLE_DOT {
			return nil, fmt.Errorf("scanSequence: found %q, expected triple dot", lit)
		}

		if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
			return nil, fmt.Errorf("scanSequence: found %q, expected close parentheses", lit)
		}
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != OF {
		p.unscan()
	} else if tok, lit = p.scanIgnoreWhitespace(); tok == IDENT {
		// NAME OF
		sequence.Of = lit
	} else if tok == BIT {
		if tok, lit = p.scanIgnoreWhitespace(); tok != STRING {
			return nil, fmt.Errorf("sequence: found %q, expected STRING", lit)
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok != GROUP_OPEN {
			p.unscan()
		} else {
			for {
				if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
					p.unscan()
					return nil, fmt.Errorf("sequence: found %q, expected IDENT", lit)
				} else {
				}

				if tok, lit := p.scanIgnoreWhitespace(); tok == PARENTHESES_OPEN {
					if tok, lit = p.scanIgnoreWhitespace(); tok != IDENT {
						p.unscan()
					} else {
					}

					if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
						return nil, fmt.Errorf("sequence: found %q, expected close parentheses", lit)
					}
				} else {
					p.unscan()
				}

				if tok, _ := p.scanIgnoreWhitespace(); tok == COMMA {
					continue
				} else {
					p.unscan()
				}

				if tok, _ := p.scanIgnoreWhitespace(); tok == GROUP_CLOSE {
					break
				} else {
					p.unscan()
				}
			}
		}

	} else if tok == INTEGER {
		sequence.Of = "INTEGER"

		if tok, _ := p.scanIgnoreWhitespace(); tok != GROUP_OPEN {
			p.unscan()
		} else {
			for {
				if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
					p.unscan()
					return nil, fmt.Errorf("sequence: found %q, expected IDENT", lit)
				} else {
				}

				if tok, lit := p.scanIgnoreWhitespace(); tok == PARENTHESES_OPEN {
					if tok, lit = p.scanIgnoreWhitespace(); tok != IDENT {
						p.unscan()
					} else {
					}

					if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
						return nil, fmt.Errorf("sequence: found %q, expected close parentheses", lit)
					}
				} else {
					p.unscan()
				}

				if tok, _ := p.scanIgnoreWhitespace(); tok == COMMA {
					continue
				} else {
					p.unscan()
				}

				if tok, _ := p.scanIgnoreWhitespace(); tok == GROUP_CLOSE {
					break
				} else {
					p.unscan()
				}
			}
		}

	} else if tok == SEQUENCE {
		sequence.Of = "SEQUENCE"
	} else if tok == CHOICE {
		sequence.Of = "CHOICE"
	} else if tok == GRAPHIC_STRING {
		sequence.Of = "GraphicString"
	} else if tok == VISIBLE_STRING {
		sequence.Of = "VisibleString"
	} else if tok == GENERALIZED_TIME {
		sequence.Of = "GeneralizedTime"
	} else if tok == NUMERIC_STRING {
		sequence.Of = "NumericString"
	} else if tok == UNIVERSAL_STRING {
		sequence.Of = "UniversalString"
	} else if tok == GENERAL_STRING {
		sequence.Of = "GeneralString"
	} else if tok == PRINTABLE_STRING {
		sequence.Of = "PrintableString"
	} else if tok == CHARACTER {
		if tok, lit = p.scanIgnoreWhitespace(); tok != STRING {
			return nil, fmt.Errorf("scanSequence: found %q, expected STRING", lit)
		}

		sequence.Of = "CHARCTER STRING"
	} else if tok == OBJECT {
		if tok, lit = p.scanIgnoreWhitespace(); tok != IDENTIFIER {
			return nil, fmt.Errorf("scanSequence: found %q, expected IDENTIFIER", lit)
		}

		sequence.Of = "OBJECT IDENTIFIER"
	} else {
		return nil, fmt.Errorf("scanSequence: found %q, expected IDENTIFIER", lit)
	}

	if err := p.scanGroup(sequence); err != nil {
		return nil, err
	}

	return sequence, nil
}

func (p *Parser) scanSet(cmmn ASNCommon) (ASNType, error) {
	// TODO: should we differentiate between ASNSet and ASNSetOf?
	set := &ASNSet{
		cmmn,
		nil,
	}

	if err := p.scanGroup(set); err != nil {
		return nil, err
	}

	return set, nil
}

func (p *Parser) scanChoice(cmmn ASNCommon) (ASNType, error) {
	choice := &ASNChoice{
		cmmn,
		nil,
	}

	if err := p.scanGroup(choice); err != nil {
		return nil, err
	}

	return choice, nil
}

func (p *Parser) scanInteger(cmmn ASNCommon) (ASNType, error) {
	obj := &ASNInteger{
		cmmn,
		ASNEnum{},
	}

	if err := p.scanEnum(obj); err != nil {
		return nil, err
	}

	// scan ranges
	// ABRT-source ::= INTEGER {service-user(0), service-provider(1)}(0..1, ...)
	if tok, _ := p.scanIgnoreWhitespace(); tok == PARENTHESES_OPEN {
		for {
			tok, lit := p.scanIgnoreWhitespace()
			if tok == TRIPLE_DOT {
				// TODO: TRIPLE_DOT
			} else if tok == IDENT {
				p.unscan()

				from, to, err := p.scanRange()
				if err != nil {
					return nil, err
				}

				_ = from
				_ = to
			} else {
				return nil, fmt.Errorf("integer: found %q, expected value", lit)
			}

			if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
				p.unscan()
				break
			}
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
			return nil, fmt.Errorf("integer: found %q, expected close parentheses", lit)
		}
	} else {
		p.unscan()
	}

	return obj, nil
}

func (p *Parser) scanEnumerated(cmmn ASNCommon) (ASNType, error) {
	obj := &ASNEnumerated{
		cmmn,
		ASNEnum{},
	}

	if err := p.scanEnum(obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (p *Parser) scanRange() (from string, to string, err error) {
	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return "", "", fmt.Errorf("range: found %q, expected value", lit)
	} else {
		from = lit
	}

	// double dot
	if tok, lit := p.scanIgnoreWhitespace(); tok != DOUBLE_DOT {
		return "", "", fmt.Errorf("range: found %q, expected double dot", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return "", "", fmt.Errorf("range: found %q, expected value", lit)
	} else {
		to = lit
	}

	return
}

func (p *Parser) scanBitString(cmmn ASNCommon) (ASNType, error) {
	if tok, lit := p.scanIgnoreWhitespace(); tok != STRING {
		p.unscan()
		return nil, fmt.Errorf("bitstring: found %q, expected STRING", lit)
	}

	obj := &ASNBitString{
		cmmn,
		ASNEnum{},
	}

	// BIT TYPE (STRING)
	if tok, lit := p.scanIgnoreWhitespace(); tok == PARENTHESES_OPEN {
		if tok, lit = p.scanIgnoreWhitespace(); tok != SIZE {
			p.unscan()
		} else if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_OPEN {
			p.unscan()
		} else {
			from, to, err := p.scanRange()
			if err != nil {
				return nil, err
			}

			_ = from
			_ = to
			if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
				return nil, fmt.Errorf("bitstring: found %q, expected close parentheses", lit)
			}
		}

		if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
			return nil, fmt.Errorf("bitstring: found %q, expected close parentheses", lit)
		}
	} else {
		p.unscan()
	}

	if tok, _ := p.scanIgnoreWhitespace(); tok == DEFAULT {
		// TODO: implement DEFAULT
	} else {
		p.unscan()
	}

	if err := p.scanEnum(obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (p *Parser) scanType(cmmn ASNCommon) (ASNType, error) {
	switch tok, lit := p.scanIgnoreWhitespace(); tok {
	case OBJECT_DESCRIPTOR:
		return &ASNObjectDescriptor{
			cmmn,
		}, nil
	case GENERAL_STRING:
		return &ASNGeneralString{
			cmmn,
		}, nil
	case NUMERIC_STRING:
		return &ASNNumericString{
			cmmn,
		}, nil
	case GRAPHIC_STRING:
		return &ASNGraphicString{
			cmmn,
		}, nil
	case VISIBLE_STRING:
		return &ASNVisibleString{
			cmmn,
		}, nil
	case UTC_TIME:
		return &ASNUTCTime{
			cmmn,
		}, nil
	case GENERALIZED_TIME:
		return &ASNGeneralizedTime{
			cmmn,
		}, nil
	case T61_STRING:
		return &ASNT61String{
			cmmn,
		}, nil
	case PRINTABLE_STRING:
		return &ASNPrintableString{
			cmmn,
		}, nil

	case OCTET:
		if tok, lit = p.scanIgnoreWhitespace(); tok != STRING {
			return nil, fmt.Errorf("type: found %q, expected IDENTIFIER", lit)
		}

		return &ASNOctetString{
			cmmn,
		}, nil
	case ENUMERATED:
		return p.scanEnumerated(cmmn)
	case SEQUENCE:
		return p.scanSequence(cmmn)
	case SET:
		return p.scanSet(cmmn)
	case CHOICE:
		return p.scanChoice(cmmn)
	case INTEGER:
		return p.scanInteger(cmmn)
	case BIT:
		return p.scanBitString(cmmn)
	case OBJECT:
		if tok, lit = p.scanIgnoreWhitespace(); tok != IDENTIFIER {
			return nil, fmt.Errorf("type: found %q, expected IDENTIFIER", lit)
		}

		return &ASNObjectIdentifier{
			cmmn,
		}, nil

	case IDENT:
		return &ASNCustom{
			cmmn,
			lit,
		}, nil
	default:
		return nil, fmt.Errorf("type: found %q, expected field", lit)
	}
}

func (p *Parser) scanDefault(e ASNEnumerer /* ASNDefaulter */) error {
	if tok, _ := p.scanIgnoreWhitespace(); tok != DEFAULT {
		p.unscan()
		return nil
	}

	// TODO: implement default
	return fmt.Errorf("Default has not been implemented yet.")
}

func (p *Parser) scanEnum(e ASNEnumerer) error {
	if tok, _ := p.scanIgnoreWhitespace(); tok != GROUP_OPEN {
		p.unscan()
		return nil
	}

	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == TRIPLE_DOT {
			// TODO: triple dot
		} else if tok == IDENT {
		} else {
			return fmt.Errorf("enum: found %q, expected IDENT %+#v", tok, lit)
		}

		name := lit

		if tok, lit = p.scanIgnoreWhitespace(); tok == PARENTHESES_OPEN {
			// CONSTANT VALUE
			if tok, lit = p.scanIgnoreWhitespace(); tok == IDENT {
			} else {
				return fmt.Errorf("enum: found %q, expected IDENT5", tok)
			}

			e.Add(name, lit)

			// CONSTANT
			if tok, lit = p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
				return fmt.Errorf("enum: found %q, expected )", tok)
			}

			if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
				p.unscan()
				break
			}

			continue
		} else {
			p.unscan()
		}

		if tok, lit = p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != GROUP_CLOSE {
		return fmt.Errorf("enum: found %q, expected GROUP_CLOSE", lit)
	}

	return nil
}

func (p *Parser) scanGroupItem(current ASNType) error {
	tok, lit := p.scanIgnoreWhitespace()

	name := lit

	position := ""

	application := false

	if tok, lit := p.scanIgnoreWhitespace(); tok == OPTIONAL_TERM_OPEN {
		for {
			tok, lit = p.scanIgnoreWhitespace()
			if tok == OPTIONAL_TERM_CLOSE {
				break
			}

			if tok == APPLICATION {
				application = true
				continue
			}

			position = lit
		}
	} else {
		p.unscan()
	}

	implicit := false
	if tok, lit = p.scanIgnoreWhitespace(); tok == IMPLICIT {
		// IMPLICIT
		implicit = true
	} else {
		p.unscan()
	}

	explicit := false
	if tok, lit = p.scanIgnoreWhitespace(); tok == EXPLICIT {
		// EXPLICIT
		explicit = true
	} else {
		p.unscan()
	}

	// TODO: EXPLICIT
	_ = explicit

	type_, err := p.scanType(ASNCommon{
		name: name,
	})
	if err != nil {
		return err
	}

	optional := false
	defaultValue := ""

	if tok, lit = p.scanIgnoreWhitespace(); tok == OPTIONAL {
		// value is optional
		optional = true
	} else if tok == DEFAULT {
		defaultValue = "DEFAULT"

		// value has a default value
		if tok, lit = p.scanIgnoreWhitespace(); tok == TRUE {
			// TODO: implement DEFAULT
		} else if tok == FALSE {
			// TODO: implement DEFAULT
		} else if tok == IDENT {
			// TODO: implement DEFAULT
		} else if tok == GROUP_OPEN {
			if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
				p.unscan()
			} else {
				// TODO: implement DEFAULT
				_ = lit
			}

			if tok, lit := p.scanIgnoreWhitespace(); tok != GROUP_CLOSE {
				return fmt.Errorf("group: found %q, expected GROUP_CLOSE", lit)
			}
		} else {
			return fmt.Errorf("group: found %q, expected VALUE", lit)
		}
	} else {
		p.unscan()
	}

	item := ASNItem{
		Name:        name,
		Position:    position,
		Optional:    optional,
		Implicit:    implicit,
		Explicit:    explicit,
		Application: application,
		Type:        type_,
		Default:     defaultValue,
	}

	switch v := current.(type) {
	case *ASNSequence:
		v.Items = append(v.Items, item)
	case *ASNChoice:
		v.Items = append(v.Items, item)
	case *ASNSet:
		v.Items = append(v.Items, item)
	default:
		return fmt.Errorf("Could not add item to type other %#v\n", current)
	}

	if tok, _ := p.scanIgnoreWhitespace(); tok != PARENTHESES_OPEN {
		p.unscan()
	} else if tok, _ := p.scanIgnoreWhitespace(); tok == CONSTRAINED {
		if tok, _ := p.scanIgnoreWhitespace(); tok != BY {
			return fmt.Errorf("group: found %q, expected BY", lit)
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok != GROUP_OPEN {
			return fmt.Errorf("group: found %q, expected GROUP_OPEN", lit)
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok != GROUP_CLOSE {
			return fmt.Errorf("group: found %q, expected GROUP_CLOSE", lit)
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok != PARENTHESES_CLOSE {
			return fmt.Errorf("group: found %q, expected PARENTHESES_CLOSE", lit)
		}
	}

	return nil
}

func (p *Parser) scanGroup(current ASNType) error {
	if tok, _ := p.scanIgnoreWhitespace(); tok != GROUP_OPEN {
		p.unscan()
		return nil
	}

	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == TRIPLE_DOT {
			item := ASNItem{
				TripleDot: true,
			}

			switch v := current.(type) {
			case *ASNSequence:
				v.Items = append(v.Items, item)
			case *ASNChoice:
				v.Items = append(v.Items, item)
			case *ASNSet:
				v.Items = append(v.Items, item)
			default:
				return fmt.Errorf("Triple dots are not implemented for type %#v\n", current)
			}
		} else if tok == IDENT {
			p.unscan()

			if err := p.scanGroupItem(current); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("group: found %q, expected IDENT %#+v", tok, lit)
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}

		// REPEATER
		if tok, lit = p.scanIgnoreWhitespace(); tok != TRIPLE_DOT {
			p.unscan()
			continue
		}

		if tok, lit = p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != GROUP_CLOSE {
		return fmt.Errorf("group: found %q, expected GROUP_CLOSE", lit)
	}

	return nil
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

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	for {
		tok, lit = p.scan()
		if tok == COMMENT {
			continue
		} else if tok == WS {
			continue
		}

		return
	}
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
