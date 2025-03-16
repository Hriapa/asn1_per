package asn1_per

import (
	"reflect"
	"testing"
)

func TestIntegerConstrain(t *testing.T) {
	var err error
	type result struct {
		intValue int
		data     []byte
		shift    uint8
	}
	type intParam struct {
		lb     int
		ub     int
		allign bool
	}
	for _, test := range []struct {
		name  string
		param intParam
		input []byte
		shift uint8
		want  result
	}{
		{
			name: `Test_Bits`,
			param: intParam{
				lb:     0,
				ub:     15,
				allign: true,
			},
			input: []byte{0x05, 0x00, 0x07},
			shift: 4,
			want: result{
				intValue: 5,
				data:     []byte{0x00, 0x07},
				shift:    0,
			},
		},
		{
			name: `Test_Bits_1`,
			param: intParam{
				lb:     0,
				ub:     15,
				allign: true,
			},
			input: []byte{0x05, 0x00, 0x07},
			shift: 2,
			want: result{
				intValue: 1,
				data:     []byte{0x05, 0x00, 0x07},
				shift:    6,
			},
		},
		{
			name: `Test_Octet`,
			param: intParam{
				lb:     0,
				ub:     255,
				allign: true,
			},
			input: []byte{0x80, 0x0b, 0x9a, 0x20},
			shift: 4,
			want: result{
				intValue: 11,
				data:     []byte{0x9a, 0x20},
				shift:    0,
			},
		},
		{
			name: `Test_Octets`,
			param: intParam{
				lb:     0,
				ub:     16777215,
				allign: true,
			},
			input: []byte{0x80, 0x0b, 0x9a, 0x20},
			shift: 0,
			want: result{
				intValue: 760352,
				data:     []byte{},
				shift:    0,
			},
		},
	} {
		res := result{}
		val := NewConstrainedInteger(test.param.lb, test.param.ub, test.param.allign)
		res.data, res.shift, err = val.Decode(test.input, test.shift)
		if err != nil {
			t.Error("error decode constrain integer")
		}
		res.intValue = val.Value
		if !reflect.DeepEqual(test.want, res) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.want, res)
			t.Fail()
		}
	}
}
