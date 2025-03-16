package asn1_per

import (
	"encoding/binary"
	"math"
	"math/big"
)

//11.5 Encoding of a constrained whole number

func constrainedWholeNumber(data []byte, shift uint8, rang int, aligned bool) (result int, outData []byte, outShift uint8, err error) {
	if len(data) == 0 {
		err = ErrorBufferToShort
		return
	}
	if !aligned {
		return unalignedConstrainedWholeNumber(data, shift, rang)
	}
	if rang > 255 {
		if shift != 0 {
			outShift = 0
			data = data[1:]
		}
	}
	// the bit-field case
	if rang < 256 {
		result, outData, outShift, err = bitFieldCase(data, shift, rang)
		return
	}
	// the one-octet case
	if rang == 256 {
		result, outData, err = oneOctetCase(data)
		return
	}
	// the two-octet case
	if 256 < rang && rang < 65536 {
		result, outData, err = twoOctetCase(data)
		return
	}
	// the indefinite length case
	result, outData, err = indefiniteLengthCase(data)
	return
}

// UNALIGNED variant

func unalignedConstrainedWholeNumber(data []byte, shift uint8, rang int) (result int, outData []byte, outShift uint8, err error) {
	bitSize := sizeInBits(rang)
	resSize := bitSize / 8
	if bitSize%8 != 0 {
		resSize += 1
	}
	byteSize := resSize
	if (shift+uint8(bitSize%8) > 8) || (shift != 0 && bitSize%8 == 0) {
		byteSize += 1
	}
	if len(data) < byteSize {
		err = ErrorBufferToShort
		return
	}
	if resSize > 8 {
		// To big Value
		err = ErrorIncorrectDecode
		return
	}
	bufRes := make([]byte, 0, resSize)
	padding := bitStringDecode(data, shift, uint8(bitSize%8), byteSize, &bufRes)

	result = int(big.NewInt(0).SetBytes(bufRes).Uint64())
	if padding == 0 {
		outData = data[byteSize:]
		outShift = 0
	} else {
		outData = data[byteSize-1:]
		outShift = 8 - padding
	}
	return
}

// the bit-field case
func bitFieldCase(data []byte, shift uint8, rang int) (result int, outData []byte, outShift uint8, err error) {
	var padding uint8
	sizeLen := uint8(sizeInBits(rang))
	if (shift + sizeLen) < 8 {
		padding = 8 - (shift + sizeLen)
		result = int((data[0] >> padding) & (0xff >> (8 - sizeLen)))
		outData = data
		outShift = shift + sizeLen
	} else {
		if len(data) < 2 {
			err = ErrorBufferToShort
			return
		}
		result = int((data[0]<<shift | data[1]>>(8-shift)) >> (8 - sizeLen))
		outData = data[1:]
		outShift = shift + sizeLen - 8
	}
	return
}

// the one-octet case
func oneOctetCase(data []byte) (result int, outData []byte, err error) {
	if len(data) == 0 {
		err = ErrorBufferToShort
		return
	}
	result = int(data[0])
	outData = data[1:]
	return
}

// the two-octet case
func twoOctetCase(data []byte) (result int, outData []byte, err error) {
	if len(data) < 2 {
		err = ErrorBufferToShort
		return
	}
	result = int(binary.BigEndian.Uint16(data[:2]))
	outData = data[2:]
	return
}

// the indefinite length case

func indefiniteLengthCase(data []byte) (result int, outData []byte, err error) {
	if len(data) < 2 {
		err = ErrorBufferToShort
		return
	}
	if data[0]>>7 == 1 {
		// Length must be less than 128 bits (16 byte). It's too big.
		err = ErrorIncorrectDecode
		return
	}
	length := data[0]
	data = data[1:]
	if length > 8 {
		// To big Value
		err = ErrorIncorrectDecode
		return
	}
	if len(data) < int(length) {
		err = ErrorBufferToShort
		return
	}
	result = int(big.NewInt(0).SetBytes(data[:length]).Uint64())
	outData = data[length:]
	return
}

func sizeInBits(rang int) int {
	x := math.Log2(float64(rang))
	return int(math.Ceil(x))
}

//11.9 General rules for encoding a length determinant

func DecodeLengthDeterminant(data []byte) (fullLength int, payload []byte, err error) {
	if len(data) < 2 {
		err = ErrorBufferToShort
		return
	}
	var length int
	if data[0]>>7 == 0 {
		length = int(data[0])
		if len(data)-1 < int(length) || length == 0 {
			err = ErrorBufferToShort
			return
		}
		fullLength = length + 1
		payload = data[1 : 1+length]
		return
	}
	if data[0]>>6 == 2 {
		length = int(binary.BigEndian.Uint16([]byte{data[0] & 0b0011_1111, data[1]}))
		if len(data)-2 < int(length) || length == 0 {
			err = ErrorBufferToShort
			return
		}
		fullLength = length + 2
		payload = data[2 : 2+length]
		return
	}
	for {
		if data[0]>>7 == 0 {
			length = int(data[0])
			if len(data)-1 < length || length == 0 {
				fullLength = 0
				payload = nil
				err = ErrorBufferToShort
				return
			}
			payload = append(payload, data[1:1+length]...)
			fullLength = fullLength + 1 + length
			break
		} else {
			length = int(binary.BigEndian.Uint16([]byte{data[0] & 0b0011_1111, data[1]})) * K16
			if len(data)-2 < length || length == 0 {
				fullLength = 0
				payload = nil
				err = ErrorBufferToShort
				return
			}
			payload = append(payload, data[2:2+length]...)
			fullLength = fullLength + 2 + length
			if len(data[2+length:]) == 0 {
				break
			}
			data = data[2+length:]
		}
	}
	return
}

// function calculate length (return number of information elemrnts in data stream and error).
// If need length in bits, bits = true
func lengthDeterminanteCalculate(data []byte, bits bool) (length int, err error) {
	if len(data) == 0 {
		err = ErrorInputParameters
		return
	}
	// length size
	var (
		lengthType     uint8
		fragmentLength int
	)
	lengthType = data[0] >> 7
	// 0xxx_xxxx - length size = 1 octet (rest 7 bits)
	if lengthType == 0 {
		length = int(data[0])
		if length == 0 {
			err = ErrorIncorrectDecode
		}
		return
	}
	lengthType = data[0] >> 6
	// 10xx_xxxx xxxx_xxxx - length size = 2 octet (rest 14 bits)
	if lengthType == 2 {
		length = int(binary.BigEndian.Uint16([]byte{data[0] & 0b0011_1111, data[1]}))
		if length == 0 {
			err = ErrorIncorrectDecode
		}
		return
	}
	byteDivisor := 1
	if bits {
		byteDivisor = 8
	}
	// 11xx_xxxx - fragmentation case
	for {
		lengthType = data[0] >> 7
		if lengthType == 0 {
			length += int(data[0])
			break
		} else {
			if lengthType != 3 {
				err = ErrorIncorrectLength
				return
			}
			if len(data) < 2 {
				err = ErrorInputParameters
				return
			}
			fragmentLength = int(binary.BigEndian.Uint16([]byte{data[0] & 0b0011_1111, data[1]})) * K16
			if fragmentLength == 0 || fragmentLength > 4 {
				err = ErrorIncorrectLength
				return
			}
			fragmentLength = fragmentLength / byteDivisor
			if len(data)-2 < fragmentLength/8 {
				err = ErrorIncorrectLength
				return
			}
			if len(data[2+fragmentLength:]) == 0 {
				break
			}
			data = data[2+fragmentLength:]
			length += fragmentLength
		}
	}
	return
}
