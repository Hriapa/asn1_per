package asn1_per

// BIT STRING with fixed length

type FixedBitString struct {
	Size     int
	Alligned bool
	Value    []uint8
}

func NewFixedBitString(size int, alligned bool) *FixedBitString {
	resSize := size / 8
	if size%8 != 0 {
		resSize += 1
	}
	return &FixedBitString{
		Size:     size,
		Alligned: alligned,
		Value:    make([]uint8, 0, resSize),
	}
}

func (b *FixedBitString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error) {
	if len(data) == 0 {
		err = ErrorBufferToShort
		return
	}
	if len(b.Value) != 0 {
		b.Value = b.Value[:0]
	}
	return fixedBitStringDecode(data, shift, b.Size, b.Alligned, &b.Value)
}

// BIT STRING with constrained length

type ConstrainedBitString struct {
	LowerBand int
	UpperBand int
	Alligned  bool
	Size      int // Result Size in Bits
	Value     []uint8
}

func NewConstrainedBitString(lb int, ub int, alligned bool) *ConstrainedBitString {
	maxSize := (ub - lb + 1) / 8
	if (ub-lb+1)%8 != 0 {
		maxSize += 1
	}
	if maxSize < 0 {
		maxSize = 1
	}
	return &ConstrainedBitString{
		LowerBand: lb,
		UpperBand: ub,
		Alligned:  alligned,
		Value:     make([]uint8, 0, maxSize),
	}
}

func (b *ConstrainedBitString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error) {
	if b.UpperBand < b.LowerBand {
		err = ErrorInputParameters
		return
	}
	if len(data) == 0 {
		err = ErrorBufferToShort
		return
	}
	if len(b.Value) != 0 {
		b.Value = b.Value[:0]
	}
	if b.UpperBand == b.LowerBand {
		return fixedBitStringDecode(data, shift, b.UpperBand, b.Alligned, &b.Value)
	}
	rang := b.UpperBand - b.LowerBand + 1

	b.Size, data, shift, err = constrainedWholeNumber(data, shift, rang, b.Alligned)
	if err != nil {
		return
	}
	b.Size += b.LowerBand
	return fixedBitStringDecode(data, shift, b.Size, b.Alligned, &b.Value)
}

// BIT STRING with unconstrained length

type UnconstrainedBitString struct {
	Alligned bool
	Size     int // Result Size in Bits
	Value    []byte
}

func NewUnconstrainedBitString(alligned bool) *UnconstrainedBitString {
	return &UnconstrainedBitString{
		Alligned: alligned,
	}
}

func (b *UnconstrainedBitString) Decode(data []byte, shift uint8) (outData []byte, outShift uint8, err error) {
	if b.Alligned {
		if shift != 0 {
			data = data[1:]
			return b.unconstrainedBitStringAlign(data)
		}
	}
	return
}

// aligned variant
func (b *UnconstrainedBitString) unconstrainedBitStringAlign(data []byte) (outData []byte, outShift uint8, err error) {
	var padding int
	b.Size, err = lengthDeterminanteCalculate(data, true)
	if err != nil {
		return
	}
	byteLength := b.Size / 8
	// padding bits for allign to octet
	padding = 8 - b.Size%8
	if padding == 8 {
		padding = 0
	}
	if padding != 0 {
		byteLength += 1
	}
	if byteLength > len(data) {
		err = ErrorIncorrectDecode
		return
	}
	b.Value = make([]byte, byteLength)

	// 1 octet length
	if b.Size < 128 {
		data = data[1:]
		return b.nonFragmentationCase(data, byteLength, uint8(padding))
	}

	// 2 octet length

	if b.Size > 127 && b.Size < 16384 {
		data = data[2:]
		return b.nonFragmentationCase(data, byteLength, uint8(padding))
	}

	// fragmentation case
	return
}

func (b *UnconstrainedBitString) nonFragmentationCase(data []byte, length int, padding uint8) (outData []byte, outShift uint8, err error) {
	if length > len(data) {
		err = ErrorIncorrectDecode
		return
	}
	if padding != 0 {
		b.allign(data, uint8(padding))
		outData = data[length-1:]
		outShift = 8 - uint8(padding)
	} else {
		copy(b.Value, data[:length])
		outData = data[length:]
		outShift = 0
	}
	return
}

func (b *UnconstrainedBitString) allign(data []byte, shift uint8) {
	b.Value[0] = data[0] >> shift
	for i := 1; i < len(b.Value); i++ {
		b.Value[i] = (data[i-1] << (8 - shift)) | data[i]>>shift
	}
}

// Декодирует BitString фиксированной длины, возвращает остаток данных и битовый сдвиг, для дальнейшего декодирования
// Необходимо указать размер в битах, трнебует ли выравнивания (Aligned PER) и указатель на результирующие данные
func fixedBitStringDecode(data []byte, shift uint8, size int, alligned bool, value *[]byte) (outData []byte, outShift uint8, err error) {
	if size > 16 && alligned {
		if shift != 0 {
			data = data[1:]
			shift = 0
		}
	}
	byteSize := fixedByteSizeCalculate(shift, size/8, uint8(size%8))
	if len(data) < byteSize {
		err = ErrorBufferToShort
		return
	}
	padding := bitStringDecode(data, shift, uint8(size%8), byteSize, value)
	if padding == 0 {
		outData = data[byteSize:]
		outShift = 0
	} else {
		outData = data[byteSize-1:]
		outShift = 8 - padding
	}
	return
}

// Принимает на вход данные, битовый сдвиг от начала байта, и остаток от деления на 8
// Заполняет value значением bitString ивозвращает количество padding bit
func bitStringDecode(data []byte, shift uint8, lastBits uint8, size int, value *[]byte) (padding uint8) {
	// количество бит в конце (p), не относящиеся к типу (до кратности байту) s s s 1 1 1 1 1 | 1 1 p p p p p p
	if (lastBits + shift) < 8 {
		padding = 8 - (lastBits + shift)
	} else {
		padding = 16 - (lastBits + shift)
	}
	if padding == 8 {
		padding = 0
	}
	// Количество нулевых бит в начале (number of padding bits in the begin of result)
	firstNulls := shift + padding
	if firstNulls > 7 {
		firstNulls -= 8
	}
	if firstNulls == 0 {
		*value = append(*value, data[:size]...)
		return
	}
	// Требуемый сдвиг от первого бита, что бы последний бит был на стыке октетов
	if shift > firstNulls {
		//Сдвиг влево
		for i := 0; i < (size - 1); i++ {
			*value = append(*value, data[i]<<(shift-firstNulls)|data[i+1]>>(8-shift-firstNulls))
		}
		(*value)[0] = (*value)[0] & (0xff >> firstNulls)

	} else {
		//Сдвиг вправо
		*value = append(*value, (data[0]<<shift)>>firstNulls)
		for i := 0; i < (size - 1); i++ {
			*value = append(*value, data[i]<<(8-padding)|data[i+1]>>(padding))
		}
	}
	return
}

// Calculate size in byte for bitstring fixed
func fixedByteSizeCalculate(shift uint8, size int, outShift uint8) int {
	if shift != 0 || outShift != 0 {
		size += 1
	}
	if (shift + outShift) > 8 {
		size += 1
	}
	return size
}

// Calculate size in byte for bitstring constrained
// Return size (int), outData (after length field) and bit shift
