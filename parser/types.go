package asn1parser

import asn1 "github.com/dutchsec/asn1"

// ASNDefinition contains the all definition types of the parsed scheme
type ASNDefinition struct {
	Name string

	Types   []ASNType
	Imports map[string][]string
}

// ASNItem is the base struct for definition types
type ASNItem struct {
	Name     string
	Position string

	Implicit    bool
	Explicit    bool
	Application bool
	Optional    bool

	Type ASNType

	Default interface{}

	TripleDot bool
}

type ASNCommon struct {
	name string

	Implicit bool
	Explicit bool

	tag asn1.ASNTag
}

func (c *ASNCommon) Name() string {
	return c.name
}

func (c *ASNCommon) Tag() asn1.ASNTag {
	return c.tag
}

type ASNEnumerer interface {
	Add(key string, v interface{})
}

type ASNType interface {
	Name() string
	Tag() asn1.ASNTag
}

type ASNAlias struct {
	ASNCommon
	Alias   string
	Default string
}

type ASNCustom struct {
	ASNCommon
	Type string
}

type ASNInteger struct {
	ASNCommon
	ASNEnum
}

type ASNEnum struct {
	Values map[string]interface{}
}

func (b *ASNEnum) Add(key string, v interface{}) {
	if b.Values == nil {
		b.Values = map[string]interface{}{}
	}

	b.Values[key] = v
}

type ASNChoice struct {
	ASNCommon

	Items []ASNItem
}

type ASNSet struct {
	ASNCommon

	Items []ASNItem
}

type ASNEnumerated struct {
	ASNCommon

	ASNEnum
}

type ASNBitString struct {
	ASNCommon

	ASNEnum
}

type ASNObjectIdentifier struct {
	ASNCommon
}

type ASNGraphicString struct {
	ASNCommon
}

type ASNT61String struct {
	ASNCommon
}

type ASNGeneralizedTime struct {
	ASNCommon
}

type ASNUTCTime struct {
	ASNCommon
}

type ASNObjectDescriptor struct {
	ASNCommon
}

type ASNPrintableString struct {
	ASNCommon
}

type ASNNumericString struct {
	ASNCommon
}

type ASNGeneralString struct {
	ASNCommon
}

type ASNVisibleString struct {
	ASNCommon
}

type ASNOctetString struct {
	ASNCommon
}

type ASNSequence struct {
	ASNCommon

	Of string

	Items []ASNItem
}
