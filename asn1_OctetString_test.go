package asn1_per

import (
	"reflect"
	"testing"
)

func TestFixedOctetString(t *testing.T) {
	type result struct {
		octetString []byte
		data        []byte
		shift       uint8
	}
	type parameters struct {
		size     int
		alligned bool
	}

	for _, test := range []struct {
		name  string
		param parameters
		input []byte
		shift uint8
		want  result
	}{
		{
			name: "Tests_1",
			param: parameters{
				size:     2,
				alligned: true,
			},
			input: []byte{0x51, 0x80, 0x07, 0x08},
			shift: 4,
			want: result{
				octetString: []byte{0x18, 0x00},
				data:        []byte{0x07, 0x08},
				shift:       4,
			},
		},
	} {
		var err error
		res := result{}
		octet := NewFixedOctetString(test.param.size, test.param.alligned)
		res.data, res.shift, err = octet.Decode(test.input, test.shift)
		if err != nil {
			t.Errorf("error decode constrain octet string")
		}
		res.octetString = octet.Value
		if !reflect.DeepEqual(test.want, res) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.want, res)
			t.Fail()
		}
	}
}

func TestConstrainedOctetString(t *testing.T) {
	type result struct {
		octetString []byte
		data        []byte
		shift       uint8
	}
	type parameters struct {
		ub       int
		lb       int
		alligned bool
	}

	for _, test := range []struct {
		name  string
		param parameters
		input []byte
		shift uint8
		want  result
	}{
		{
			name: "Tests_1",
			param: parameters{
				ub:       8,
				lb:       3,
				alligned: true,
			},
			input: []byte{0xa0, 0xaf, 0x20, 0x60, 0x52, 0xf0, 0x99, 0x03, 0xb3},
			shift: 0,
			want: result{
				octetString: []byte{0xaf, 0x20, 0x60, 0x52, 0xf0, 0x99, 0x03, 0xb3},
				data:        []byte{},
				shift:       0,
			},
		},
		{
			name: "Tests_1",
			param: parameters{
				ub:       8,
				lb:       3,
				alligned: true,
			},
			input: []byte{0x10, 0xaf, 0x20, 0x60, 0x52, 0xf0, 0x99, 0x03, 0xb3},
			shift: 2,
			want: result{
				octetString: []byte{0xaf, 0x20, 0x60, 0x52, 0xf0},
				data:        []byte{0x99, 0x03, 0xb3},
				shift:       0,
			},
		},
	} {
		var err error
		res := result{}
		octet := NewConstrainedOctetString(test.param.lb, test.param.ub, test.param.alligned)
		res.data, res.shift, err = octet.Decode(test.input, test.shift)
		if err != nil {
			t.Errorf("error decode constrain octet string")
		}
		res.octetString = octet.Value
		if !reflect.DeepEqual(test.want, res) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.want, res)
			t.Fail()
		}
	}
}
