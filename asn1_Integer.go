package asn1_per

// INTEGER Type Сonstrained

type ConstrainedInteger struct {
	LowerBand int
	UpperBand int
	Alligned  bool
	Value     int
}

func NewConstrainedInteger(lb int, ub int, alligned bool) *ConstrainedInteger {
	return &ConstrainedInteger{
		LowerBand: lb,
		UpperBand: ub,
		Alligned:  alligned,
		Value:     0,
	}
}

func (c *ConstrainedInteger) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error) {
	if c.LowerBand >= c.UpperBand {
		err = ErrorInputParameters
		return
	}
	if len(data) == 0 {
		err = ErrorBufferToShort
		return
	}

	rang := c.UpperBand - c.LowerBand + 1

	if !c.Alligned {
		c.Value, outData, outShift, err = unalignedConstrainedWholeNumber(data, shift, rang)
		return
	}

	// the bit filed case
	if rang < 256 {
		c.Value, outData, outShift, err = bitFieldCase(data, shift, rang)
		return
	}
	// the one octet case
	if rang == 256 {
		if shift != 0 {
			outShift = 0
			data = data[1:]
		}
		c.Value, outData, err = oneOctetCase(data)
		return
	}
	// the two-octet case
	if rang > 256 && rang < 65536 {
		if shift != 0 {
			outShift = 0
			data = data[1:]
		}
		c.Value, outData, err = twoOctetCase(data)
		return
	}
	// other variants
	var lengthSize uint8
	lengthSize, err = lengthSizeCalculate(rang)
	if err != nil {
		return
	}
	// length mo than 256 octets is to big and has not of any practical interest
	// in that case lenght is not alligned of octets
	if len(data) < 2 {
		err = ErrorBufferToShort
		return
	}
	// low band of integer length = 1
	size := ((data[0]<<shift | data[1]>>(8-shift)) >> (8 - lengthSize)) + 1
	data = data[1:]
	shift += lengthSize
	if shift > 8 {
		data = data[1:]
	}
	if len(data) < int(size) {
		err = ErrorIncorrectLength
		return
	}
	c.Value = int(parseInt64(data[:size]))
	outData = data[size:]
	return
}

// return number of bits nedded for length (number of bytes) coding
func lengthSizeCalculate(rang int) (size uint8, err error) {
	// bytes neds for coding integer value
	sizeinBits := sizeInBits(rang)
	sizeinBytes := sizeinBits / 8
	if sizeinBits%8 != 0 {
		sizeinBytes += 1
	}
	// not of any practical interest
	if sizeinBytes > 127 {
		err = ErrorInputParameters
		return
	}
	size = uint8(sizeInBits(sizeinBytes))
	return
}

func parseInt64(in []byte) (val int64) {
	for i := 0; i < len(in); i++ {
		val <<= 8
		val |= int64(in[i])
	}
	// Shift up and down in order to sign extend the result.
	val <<= 64 - uint8(len(in))*8
	val >>= 64 - uint8(len(in))*8
	return
}
