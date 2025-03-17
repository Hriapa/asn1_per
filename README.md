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
```go
   data := []byte{0x02, 0x1,0x02,0x3,0x04}, 
   
   fullLength, payload, err := DecodeLengthDeterminant(data)
   if err != nil{
    //processing error
   }
```
```
Result:   
   fullLength = 2
   payload = []byte{0x1,0x02}
   err = nil
```

## Decode Functions Parameters

All decode functions (for different types) have the same input and output parameters

```go
    Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error)
```
data - input raw data  
shift - bit shift (number of end bits from type placed before)  
outData - rest data   
outShift - bit shift (number of end bits from INTEGER type (0 - if payload end on the corner of octet))  

```
Example:
    coding value = []byte{0b1111_1111,0b1111_1111}
    previos_data = 0b0000

    data = []byte{0b0000_1111,0b1111_1111,0b1111_0000}
    shift = 4
    outData = []byte{0b1111_0000}
    outshift = 4
```


## Types

### INTEGER

#### ConstrainedInteger

```go
    func NewConstrainedInteger(lb int, ub int, alligned bool) *ConstrainedInteger
```
lb - lower band  
ub - upper band  
alligned - true:alligned format\false:not alligned format  
Include Value field with integer type  

Decoding 

```go
    func (c *ConstrainedInteger) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error)
```

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

#### FixedBitString

BIT STRING with fixed length

```go
    func NewFixedBitString(size int, alligned bool) *FixedBitString
```
size - size in bits   
alligned - true:alligned format\false:not alligned format  
Include Value field with []uint8 type alligned by the end of octet

Decoding 

```go
    func (b *FixedBitString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error)
```

Eexample:

```
exapmple := SEQUENCE{
    INTEGER(0..15)
    BIT STRING (SIZE (5))
    OCTET STRING (SIZE (2))
}
```

```go
    data := []byte{0x51, 0x80, 0x07, 0x08}
    bit := NewFixedBitString(5, true)

    out, shift, err := bit.Decode(data, 4)
    if err!= nil{
        // error processing
    }
    value := bit.Value
```
```
Result:
    value = []uint8{0x03}
    out = []byte{0x80, 0x07, 0x08}
    shift = 1
    err = nil
```

#### ConstrainedBitString

BIT STRING with constrained length

```go
    func NewConstrainedBitString(lb int, ub int, alligned bool) *ConstrainedBitString
```
lb - lower band  
ub - upper band  
alligned - true:alligned format\false:not alligned format  
Include Value field with []uint8 type alligned by the end of octet and Bits size with integer type

Decoding 

```go
    func (b *ConstrainedBitString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error)
```

Eexample:

```
exapmple := SEQUENCE{
    INTEGER(0..1)
    BIT STRING (SIZE (1..160))
    OCTET STRING (SIZE (2))
}
```

```go
    data := []byte{0x0f, 0x80, 0x4e, 0x19, 0x69, 0x72, 0x2b, 0xa8}
    bit := NewConstrainedBitString(1, 160, true)

    out, shift, err := bit.Decode(data, 1)
    if err!= nil{
        // error processing
    }
    value := bit.Value
    size := bit.Size
```
```
Result:
    value = []uint8{0x4e, 0x19, 0x69, 0x72}
    size = 32
    out = []byte{0x2b, 0xa8}
    shift = 0
    err = nil
```

#### UnconstrainedBitString

BIT STRING with unconstrained length

```go
    func NewUnconstrainedBitString(alligned bool) *UnconstrainedBitString
```

alligned - true:alligned format\false:not alligned format  
Include Value field with []byte type alligned by the end of octet

__!!! ONLY ALLIGNED VARIANT (Non aligned, not realized yet) !!!__ 

Decoding 

```go
    func (b *UnconstrainedBitString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error)
```

Eexample:

```
exapmple := SEQUENCE{
    INTEGER(0..1)
    BIT STRING 
    OCTET STRING (SIZE (2))
}
```

```go
    data := []byte{0x0f, 0x80, 0x80, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x2b, 0xa8},
    bit := NewUnconstrainedBitString(true)

    out, shift, err := bit.Decode(data, 1)
    if err!= nil{
        // error processing
    }
    value := bit.Value
    size := bit.Size
```
```
Result:
    value = []byte{0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72}
    size = 128
    out = []byte{0x2b, 0xa8}
    shift = 0
    err = nil
```

### OCTET STRING

#### FixedOctetString

OCTET STRING with fixed length

```go
    func NewFixedOctetString(size int, alligned bool) *FixedOctetString
```
size - size in bits   
alligned - true:alligned format\false:not alligned format  
Include Value field with []byte type

Decoding 

```go
   func (o *FixedOctetString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error)
```

Eexample:

```
exapmple := SEQUENCE{
    INTEGER(0..15)
    OCTET STRING (SIZE (2))
}
```

```go
    data := []byte{0x51, 0x80, 0x07, 0x08}
    octet := NewFixedOctetString(2, true)

    out, shift, err := octet.Decode(data, 4)
    if err!= nil{
        // error processing
    }
    value := octet.Value
```
```
Result:
    value = []uint8{0x18, 0x00}
    out = []byte{0x07, 0x08}
    shift = 4
    err = nil
```

#### ConstrainedOctetString

OCTET STRING with constrained length

```go
    func NewConstrainedOctetString(lb int, ub int, alligned bool) *ConstrainedOctetString
```
lb - lower band  
ub - upper band  
alligned - true:alligned format\false:not alligned format  
Include Value field with []byte type

Decoding 

```go
   func (o *ConstrainedOctetString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error)
```

Eexample:

```
exapmple := SEQUENCE{
    INTEGER(0..3)
    OCTET STRING (SIZE (3..8))
}
```

```go
    data := []byte{0x10, 0xaf, 0x20, 0x60, 0x52, 0xf0, 0x99, 0x03, 0xb3}
    octet := NewConstrainedOctetString(3, 8, true)

    out, shift, err := octet.Decode(data, 2)
    if err!= nil{
        // error processing
    }
    value := octet.Value
```
```
Result:
    value = []uint8{0xaf, 0x20, 0x60, 0x52, 0xf0}
    out = []byte{0x99, 0x03, 0xb3}
    shift = 0
    err = nil
```