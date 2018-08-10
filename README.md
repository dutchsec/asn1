# ASN1 Scheme Parser [![](https://godoc.org/github.com/dutchsec/asn1-scheme-parser?status.svg)](http://godoc.org/github.com/dutchsec/asn1-scheme-parser) [![Go Report Card](https://goreportcard.com/badge/dutchsec/asn1-scheme-parser)](https://goreportcard.com/report/dutchsec/asn1-scheme-parser) 

The ASN1 scheme parser will parse an ASN1 file and return the definition. The definition can be used to parse ASN1 encoded data structures. 

## Usage

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

## Copyright and license

Code released under [Apache License 2.0](LICENSE).
