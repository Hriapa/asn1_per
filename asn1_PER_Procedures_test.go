package asn1_per

import (
	"reflect"
	"testing"
)

// unaligned constrained whole number
func TestUnalignedConstrainedWholeNumber(t *testing.T) {
	type testData struct {
		data  []byte
		shift uint8
		value int
	}
	var err error
	for _, test := range []struct {
		name  string
		input testData
		want  testData
	}{
		{
			name: `Test1`,
			input: testData{
				data:  []byte{0x0b, 0x22},
				shift: 3,
				value: 5,
			},
			want: testData{
				data:  []byte{0x0b, 0x22},
				shift: 6,
				value: 2,
			},
		},
		{
			name: `Test2`,
			input: testData{
				data:  []byte{0x0b, 0x22, 0x13, 0xff},
				shift: 3,
				value: 420,
			},
			want: testData{
				data:  []byte{0x22, 0x13, 0xff},
				shift: 4,
				value: 178,
			},
		},
		{
			name: `Test3`,
			input: testData{
				data:  []byte{0x0b, 0x22, 0x13, 0xff},
				shift: 6,
				value: 131071,
			},
			want: testData{
				data:  []byte{0x13, 0xff},
				shift: 7,
				value: 102665,
			},
		},
	} {
		result := testData{}
		result.value, result.data, result.shift, err = unalignedConstrainedWholeNumber(test.input.data, test.input.shift, test.input.value)
		if err != nil {
			t.Errorf("error decode unaligned constrained whole number")
		}
		if !reflect.DeepEqual(test.want, result) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.want, result)
			t.Fail()
		}
	}
}
