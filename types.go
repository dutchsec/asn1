package asn1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// BIT STRING
func NewBitString() BitString {
	return BitString{}
}

// BitString is the structure to use when you want an ASN.1 BIT STRING type. A
// bit string is padded up to the nearest byte in memory and the number of
// valid bits is recorded. Padding bits will be zero.
type BitString struct {
	Bytes     []byte // bits packed into bytes.
	BitLength int    // length in bits.
}

// At returns the bit at the given index. If the index is out of range it
// returns 0.
func (b BitString) At(i int) int {
	if i < 0 || i >= b.BitLength {
		return 0
	}
	x := i / 8
	y := 7 - uint(i%8)
	return int(b.Bytes[x]>>y) & 1
}

// RightAlign returns a slice where the padding bits are at the beginning. The
// slice may share memory with the BitString.
func (b BitString) RightAlign() []byte {
	shift := uint(8 - (b.BitLength % 8))
	if shift == 8 || len(b.Bytes) == 0 {
		return b.Bytes
	}

	a := make([]byte, len(b.Bytes))
	a[0] = b.Bytes[0] >> shift
	for i := 1; i < len(b.Bytes); i++ {
		a[i] = b.Bytes[i-1] << (8 - shift)
		a[i] |= b.Bytes[i] >> shift
	}

	return a
}

func (s *BitString) UnmarshalRawValue(rv *RawValue) error {
	data := rv.Content

	if len(data) == 0 {
		return fmt.Errorf("zero length BIT STRING")
	}

	paddingBits := int(data[0])
	if paddingBits > 7 ||
		len(data) == 1 && paddingBits > 0 ||
		data[len(data)-1]&((1<<data[0])-1) != 0 {
		return fmt.Errorf("invalid padding bits in BIT STRING")
	}

	var obj BitString
	obj.BitLength = (len(data)-1)*8 - paddingBits
	obj.Bytes = data[1:]
	*s = obj

	return nil
}

/*
func (ctx *Context) encodeBitString(value reflect.Value) ([]byte, error) {
	bitString, ok := value.Interface().(BitString)
	if !ok {
		return nil, wrongType(bitStringType.String(), value)
	}

	data := make([]byte, len(bitString.Bytes)+1)
	// As the first octet, we encode the number of unused bits at the end.
	data[0] = byte(8 - (bitString.BitLength % 8))
	copy(data[1:], bitString.Bytes)
	return data, nil
}

func (ctx *Context) decodeBitString(data []byte, value reflect.Value) error {
	// TODO check value type
	if len(data) == 0 {
		return syntaxError("zero length BIT STRING")
	}
	paddingBits := int(data[0])
	if paddingBits > 7 ||
		len(data) == 1 && paddingBits > 0 ||
		data[len(data)-1]&((1<<data[0])-1) != 0 {
		return syntaxError("invalid padding bits in BIT STRING")
	}
	var obj BitString
	obj.BitLength = (len(data)-1)*8 - paddingBits
	obj.Bytes = data[1:]

	value.Set(reflect.ValueOf(obj))
	return nil
}
*/

// Oid is used to encode and decode ASN.1 OBJECT IDENTIFIERs.
type Oid []uint

// Cmp returns zero if both Oids are the same, a negative value if oid
// lexicographically precedes other and a positive value otherwise.
func (oid Oid) Cmp(other Oid) int {
	for i, n := range oid {
		if i >= len(other) {
			return 1
		}
		if n != other[i] {
			return int(n) - int(other[i])
		}
	}
	return len(oid) - len(other)
}

// String returns the dotted representation of oid.
func (oid Oid) String() string {
	if len(oid) == 0 {
		return ""
	}
	s := fmt.Sprintf(".%d", oid[0])
	for i := 1; i < len(oid); i++ {
		s += fmt.Sprintf(".%d", oid[i])
	}
	return s
}

/*
func (ctx *Context) encodeOid(value reflect.Value) ([]byte, error) {
	// Check values
	oid, ok := value.Interface().(Oid)
	if !ok {
		return nil, wrongType(oidType.String(), value)
	}

	value1 := uint(0)
	if len(oid) >= 1 {
		value1 = oid[0]
		if value1 > 2 {
			return nil, parseError("invalid value for first element of OID: %d", value1)
		}
	}

	value2 := uint(0)
	if len(oid) >= 2 {
		value2 = oid[1]
		if value2 > 39 {
			return nil, parseError("invalid value for first element of OID: %d", value2)
		}
	}

	bytes := []byte{byte(40*value1 + value2)}
	for i := 2; i < len(oid); i++ {
		bytes = append(bytes, encodeMultiByteTag(oid[i])...)
	}
	return bytes, nil
}

func (ctx *Context) decodeOid(data []byte, value reflect.Value) error {
	// TODO check value type
	if len(data) == 0 {
		value.Set(reflect.ValueOf(Oid{}))
		return nil
	}

	value1 := uint(data[0] / 40)
	value2 := uint(data[0]) - 40*value1
	oid := Oid{value1, value2}

	reader := bytes.NewBuffer(data[1:])
	for reader.Len() > 0 {
		valueN, err := decodeMultiByteTag(reader)
		if err != nil {
			return parseError("invalid value element in Object Identifier")
		}
		oid = append(oid, valueN)
	}

	value.Set(reflect.ValueOf(oid))
	return nil
}

*/

// Null is used to encode and decode ASN.1 NULLs.
type Null struct{}

func (s *Null) UnmarshalRawValue(rv *RawValue) error {
	return nil
}

type Real struct {
	string
}

func (s *Real) UnmarshalRawValue(rv *RawValue) error {
	*s = Real{
		string(rv.Content),
	}
	return nil
}

type FloatingPoint struct {
	string
}

func (s *FloatingPoint) UnmarshalRawValue(rv *RawValue) error {
	*s = FloatingPoint{
		string(rv.Content),
	}
	return nil
}

type ANY []byte

func (s *ANY) UnmarshalRawValue(rv *RawValue) error {
	*s = rv.Content
	return nil
}

type ObjectDescriptor struct {
	string
}

func (s *ObjectDescriptor) UnmarshalRawValue(rv *RawValue) error {
	*s = ObjectDescriptor{
		string(rv.Content),
	}

	return nil
}

type PrintableString struct {
	string
}

func (s *PrintableString) UnmarshalRawValue(rv *RawValue) error {
	*s = PrintableString{
		string(rv.Content),
	}

	return nil
}

type GraphicString struct {
	string
}

func (s *GraphicString) UnmarshalRawValue(rv *RawValue) error {
	*s = GraphicString{
		string(rv.Content),
	}

	return nil
}

type GeneralString struct {
	string
}

func (s *GeneralString) UnmarshalRawValue(rv *RawValue) error {
	*s = GeneralString{
		string(rv.Content),
	}

	return nil
}

type T61String struct {
	string
}

func (s *T61String) UnmarshalRawValue(rv *RawValue) error {
	*s = T61String{
		string(rv.Content),
	}

	return nil
}

type GeneralizedTime struct {
	string
}

func (s *GeneralizedTime) UnmarshalRawValue(rv *RawValue) error {
	*s = GeneralizedTime{
		string(rv.Content),
	}

	return nil
}

type UTCTime struct {
	string
}

func (s *UTCTime) UnmarshalRawValue(rv *RawValue) error {
	*s = UTCTime{
		string(rv.Content),
	}

	return nil
}

type IA5String struct {
	string
}

func (s *IA5String) UnmarshalRawValue(rv *RawValue) error {
	*s = IA5String{
		string(rv.Content),
	}

	return nil
}

type OctetString struct {
	string
}

func (s *OctetString) UnmarshalRawValue(rv *RawValue) error {
	*s = OctetString{
		string(rv.Content),
	}
	return nil
}

func (s *OctetString) String() string {
	return s.string
}

type UTF8String struct {
	string
}

func (s *UTF8String) UnmarshalRawValue(rv *RawValue) error {
	*s = UTF8String{
		string(rv.Content),
	}
	return nil
}

type ObjectIdentifier struct {
	string
}

func (s *ObjectIdentifier) UnmarshalRawValue(rv *RawValue) error {
	if len(rv.Content) == 0 {
		return fmt.Errorf("Not enough bytes for ObjectIdentifier")
	}

	vals := make([]int, 2)
	vals[0] = int(rv.Content[0]) / 40
	vals[1] = int(rv.Content[0]) % 40

	value := 0
	for i := 1; i < len(rv.Content); i++ {
		b := int(rv.Content[i])

		value = value<<7 + b&127
		if b&128 == 128 {
			continue
		}

		vals = append(vals, value)

		// reset
		value = 0
	}

	svals := make([]string, len(vals))
	for i, _ := range vals {
		svals[i] = fmt.Sprintf("%d", vals[i])
	}

	*s = ObjectIdentifier{
		strings.Join(svals, "."),
	}

	fmt.Printf("ObjectIdentifier %s %x\n", *s, rv.Content)

	return nil
}

type VisibleString struct {
	string
}

func (s *VisibleString) UnmarshalRawValue(rv *RawValue) error {
	*s = VisibleString{
		string(rv.Content),
	}
	return nil
}

func (s *VisibleString) String() string {
	return s.string
}

type Bool struct {
	bool
}

var (
	BoolTrue  = Bool{true}
	BoolFalse = Bool{false}
)

func (s *Bool) UnmarshalRawValue(rv *RawValue) error {
	data := rv.Content

	if true /* der encoding */ {
		boolValue := parseBigInt(data).Cmp(big.NewInt(0)) != 0
		if boolValue {
			*s = BoolTrue
		} else {
			*s = BoolFalse
		}
		return nil
	}

	// DER is more restrict regarding valid booleans
	if len(data) == 1 {
		switch data[0] {
		case 0x00:
			*s = BoolFalse
			return nil
		case 0xff:
			*s = BoolTrue
			return nil
		}
	}

	return fmt.Errorf("Unexpected bool value: %d", data[0])
}

type Integer struct {
	int64
}

func (s Integer) Int64() int64 {
	return s.int64
}

func (s *Integer) UnmarshalRawValue(rv *RawValue) error {
	data := rv.Content

	if len(data) > 8 {
		return fmt.Errorf("integer too large for Go type 'int64'")
	}

	// Sign extend the value
	extensionByte := byte(0x00)
	if len(data) > 0 && data[0]&0x80 != 0 {
		extensionByte = byte(0xff)
	}

	extension := make([]byte, 8-len(data))
	for i := range extension {
		extension[i] = extensionByte
	}

	data = append(extension, data...)
	// Decode binary
	num := int64(0)
	for i := 0; i < len(data); i++ {
		num <<= 8
		num |= int64(data[i])
	}

	*s = Integer{num}
	return nil
}

var ErrUnparsedObjects = errors.New("Unparsed objects")
