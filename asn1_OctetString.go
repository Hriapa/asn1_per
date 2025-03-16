package asn1_per

// OCTET STRING with fixed length

type FixedOctetString struct {
	Size     int
	Alligned bool
	Value    []byte
}

func NewFixedOctetString(size int, alligned bool) *FixedOctetString {
	return &FixedOctetString{
		Size:     size,
		Alligned: alligned,
		Value:    make([]byte, 0, size),
	}
}

func (o *FixedOctetString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error) {
	if len(data) == 0 {
		err = ErrorBufferToShort
		return
	}
	if len(o.Value) != 0 {
		o.Value = o.Value[:0]
	}
	return fixedOctetStringDecode(data, shift, o.Size, o.Alligned, &o.Value)
}

// OCTET STRING with constrained length

type ConstrainedOctetString struct {
	LowerBand int
	UpperBand int
	Alligned  bool
	Value     []byte
}

func NewConstrainedOctetString(lb int, ub int, alligned bool) *ConstrainedOctetString {
	maxSize := ub - lb + 1
	if maxSize < 0 {
		maxSize = 1
	}
	return &ConstrainedOctetString{
		LowerBand: lb,
		UpperBand: ub,
		Alligned:  alligned,
		Value:     make([]byte, 0, maxSize),
	}
}

func (o *ConstrainedOctetString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error) {
	if o.UpperBand < o.LowerBand {
		err = ErrorInputParameters
		return
	}
	if len(data) == 0 {
		err = ErrorBufferToShort
		return
	}
	if len(o.Value) != 0 {
		o.Value = o.Value[:0]
	}
	if o.UpperBand == o.LowerBand {
		return fixedOctetStringDecode(data, shift, o.UpperBand, o.Alligned, &o.Value)
	}
	var size int
	rang := o.UpperBand - o.LowerBand + 1

	size, data, shift, err = constrainedWholeNumber(data, shift, rang, o.Alligned)
	if err != nil {
		return
	}
	size += o.LowerBand
	return fixedOctetStringDecode(data, shift, size, o.Alligned, &o.Value)
}

// Decode octed string wiht fixed length
func fixedOctetStringDecode(data []byte, shift uint8, size int, alligned bool, value *[]byte) (outData []byte, outShift uint8, err error) {
	if size > 2 && alligned {
		if shift != 0 {
			data = data[1:]
			shift = 0
		}
	}
	if len(data) < size {
		err = ErrorBufferToShort
		return
	}
	if shift == 0 {
		*value = append(*value, data[:size]...)
		outShift = 0
		outData = data[size:]
		return
	}
	for i := 0; i < size-1; i++ {
		*value = append(*value, data[i]<<(shift)|data[i+1]>>(8-shift))
	}
	outShift = shift
	outData = data[size-1:]
	return
}
