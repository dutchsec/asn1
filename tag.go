package asn1

func Tag(class ASNClass, value ASNValue) ASNTag {
	return ASNTag{
		Class: class,
		Value: value,
	}
}
