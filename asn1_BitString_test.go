package asn1_per

import (
	"reflect"
	"testing"
)

// BIT STRING with fixed length
func TestFixedBitString(t *testing.T) {
	var err error
	type bitStringParameters struct {
		size     int
		alligned bool
	}
	type result struct {
		bitString []byte
		outData   []byte
		outShift  uint8
	}
	for _, test := range []struct {
		name       string
		parameters bitStringParameters
		inputData  []byte
		inputShift uint8
		result     result
	}{
		{
			name: `Test_1`,
			parameters: bitStringParameters{
				size:     28,
				alligned: true,
			},
			inputData:  []byte{0xbf, 0xaf, 0x20, 0x60, 0x52, 0xf0, 0x99, 0x03, 0xb3},
			inputShift: 0,
			result: result{
				bitString: []byte{0x0b, 0xfa, 0xf2, 0x06},
				outData:   []byte{0x60, 0x52, 0xf0, 0x99, 0x03, 0xb3},
				outShift:  4,
			},
		},
		{
			name: `Test_2`,
			parameters: bitStringParameters{
				size:     5,
				alligned: true,
			},
			inputData:  []byte{0x85, 0x24, 0xab},
			inputShift: 5,
			result: result{
				bitString: []byte{0x14},
				outData:   []byte{0x24, 0xab},
				outShift:  2,
			},
		},
		// {
		// 	name: `Test_3`,
		// 	parameters: bitStringParameters{
		// 		size:     0,
		// 		alligned: true,
		// 	},
		// 	inputData:  []byte{},
		// 	inputShift: 0,
		// 	result: result{
		// 		bitString: []byte{},
		// 		outData:   []byte{},
		// 		outShift:  0,
		// 	},
		// },
	} {
		res := result{}
		b := NewFixedBitString(test.parameters.size, test.parameters.alligned)
		res.outData, res.outShift, err = b.Decode(test.inputData, test.inputShift)
		if err != nil {
			t.Errorf(`error fixed bit string decode`)
		}
		res.bitString = b.Value
		if !reflect.DeepEqual(test.result, res) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.result, res)
			t.Fail()
		}
	}
}

// BIT STRING with constrained length
func TestConstrainedBitString(t *testing.T) {
	var err error
	type bitStringParameters struct {
		lowerBand int
		upperBand int
		alligned  bool
	}
	type result struct {
		bitString []byte
		bitSize   int
		outData   []byte
		outShift  uint8
	}
	for _, test := range []struct {
		name       string
		parameters bitStringParameters
		inputData  []byte
		inputShift uint8
		result     result
	}{
		{
			name: `Test_1`,
			parameters: bitStringParameters{
				lowerBand: 1,
				upperBand: 160,
				alligned:  true,
			},
			inputData:  []byte{0x0f, 0x80, 0x4e, 0x19, 0x69, 0x72, 0x2b, 0xa8},
			inputShift: 1,
			result: result{
				bitString: []byte{0x4e, 0x19, 0x69, 0x72},
				bitSize:   32,
				outData:   []byte{0x2b, 0xa8},
				outShift:  0,
			},
		},
	} {
		res := result{}
		b := NewConstrainedBitString(test.parameters.lowerBand, test.parameters.upperBand, test.parameters.alligned)
		res.outData, res.outShift, err = b.Decode(test.inputData, test.inputShift)
		if err != nil {
			t.Errorf(`error fixed bit string decode`)
		}
		res.bitString = b.Value
		res.bitSize = b.Size
		if !reflect.DeepEqual(test.result, res) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.result, res)
			t.Fail()
		}
	}
}

// BIT STRING with unconstrained length

func TestUnconstrainedBitString(t *testing.T) {
	var err error
	type result struct {
		bitString []byte
		bitSize   int
		outData   []byte
		outShift  uint8
	}
	for _, test := range []struct {
		name       string
		allign     bool
		inputData  []byte
		inputShift uint8
		result     result
	}{
		{
			name:       `Test_1`,
			allign:     true,
			inputData:  []byte{0x0f, 0x20, 0x4e, 0x19, 0x69, 0x72, 0x2b, 0xa8},
			inputShift: 1,
			result: result{
				bitString: []byte{0x4e, 0x19, 0x69, 0x72},
				bitSize:   32,
				outData:   []byte{0x2b, 0xa8},
				outShift:  0,
			},
		},
		{
			name:       `Test_2`,
			allign:     true,
			inputData:  []byte{0x0f, 0x80, 0x80, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x2b, 0xa8},
			inputShift: 1,
			result: result{
				bitString: []byte{0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72, 0x4e, 0x19, 0x69, 0x72},
				bitSize:   128,
				outData:   []byte{0x2b, 0xa8},
				outShift:  0,
			},
		},
	} {
		res := result{}
		b := NewUnconstrainedBitString(test.allign)
		res.outData, res.outShift, err = b.Decode(test.inputData, test.inputShift)
		if err != nil {
			t.Errorf(`error fixed bit string decode`)
		}
		res.bitString = b.Value
		res.bitSize = b.Size
		if !reflect.DeepEqual(test.result, res) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.result, res)
			t.Fail()
		}
	}
}
