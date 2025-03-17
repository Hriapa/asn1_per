# ASN. 1 PER (Packets Encoding Rules) Decoding functions (for some Types) for the GO programming language.

## Introduction

    Decoding methods for 3 types encoded using ASN.1 PER encoding rules (aligned and unaligned format).

## Common Functions    

### Encoding Unconstrained Lenght according rules X.691 (p. 11.9)

```go
    func DecodeLengthDeterminant(data []byte) (fullLength int, payload []byte, err error)
```
Input: raw data,
Output:
    fullLength - length of payload
    payload - raw data
    err - error
Example:

   if data = []byte{0x02, 0x1,0x02,0x3,0x04}, 
   thhen
   fullLength = 2
   payload = []byte{0x1,0x02}
   err = nil

## Types

### INTEGER

#### ConstrainedInteger

```go
    func NewConstrainedInteger(lb int, ub int, alligned bool) *ConstrainedInteger
```
lb - lower band /n
ub - upper band /n
alligned - true:alligned format\false:not alligned format 

Decoding 

```go
    func (c *ConstrainedInteger) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error)
```

data - input raw data
shift - bit shift (number of end bits from type placed before)
outData - rest data 
outShift - bit shift (number of end bits from INTEGER type (0 - if payload end on the corner of octet))

Eexample:

```
exapmple := SEQUENCE{
    BIT STRING (SIZE (4))
    INTEGER(0..15)
    OCTET STRING (SIZE (2))
}
```

```go

    data := []byte{0x05, 0x00, 0x07}
    integer := NewConstrainedInteger(0, 15, true)

    out, shift, err := integer.Decode(data, 4)
    if err!= nil{
        // error processing
    }
    value := integer.Value
```
```
Result:
    value = 5
    out = []byte{0x00, 0x07}
    shift = 0
    err = nil
```

### BIT STRING