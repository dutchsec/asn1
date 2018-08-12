# ASN1 [![](https://godoc.org/github.com/dutchsec/asn1?status.svg)](http://godoc.org/github.com/dutchsec/asn1) [![Go Report Card](https://goreportcard.com/badge/dutchsec/asn1)](https://goreportcard.com/report/dutchsec/asn1) 

## ASN1 

## ASN1 Code Generator

## ASN1 Scheme Parser
The ASN1 scheme parser will parse an ASN1 file and return the definition. The definition can be used to parse ASN1 encoded data structures. 

### Usage

This fragment will parse the ASN1 scheme and return a definition.

```
r, err := os.Open(arg)
if err != nil {
    panic(err)
}

parser := asn1parser.NewParser(r)

def, err := parser.Parse()
if err != nil {
    panic(err)
}
```

## Sponsors

This project has been made possible by Sentryo and Dutchsec. 

## Contributors

* [Remco Verhoef](https://twitter.com/remco_verhoef)


Parts of the ASN1 decoding have been included from https://github.com/Logicalis/asn1.

## Copyright and license

Code released under [Apache License 2.0](LICENSE).
