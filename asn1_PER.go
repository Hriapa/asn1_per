package asn1_per

import (
	"errors"
)

// As Recommendation ITU-T X.691

var (
	ErrorBufferToShort   = errors.New("ASN.1 to short input buffer")
	ErrorIncorrectLength = errors.New("ASN.1 incorrect length")
	ErrorBigLength       = errors.New("ASN.1 length too big")
	ErrorIncorrectDecode = errors.New("ASN.1 format decode incorrect")
	ErrorShiftIncorrect  = errors.New("ASN.1 incorrect bitShiftValue")
	ErrorInputParameters = errors.New("ASN.1 incorrect input parameters")
)

const K16 = 16383
