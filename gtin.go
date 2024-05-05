/*
Package gtin implements a function to convert a string to a GTIN-14 number.

GTINs can be 8, 12, 13 or 14 digits in length and are known as GTIN-8, GTIN-12, GTIN-13 and GTIN-14.
The first digits are a GS1 Prefix assigned to a GS1 member organization. The last digit is a check digit.
The rest of the digits are the item reference and have no meaning, except to identify an item.

All GTINs can be represented with 14 digits, using zero padding.

Here are the different GTINs:

GTIN-8
- Carries EAN-8 barcodes.
- 7 digits + check digit
- First digits is GS1-8 Prefix

GTIN-12
- Carries UPC-A barcode
- 11 digits + check digit
- First digits is UPC Prefix

GTIN-13
- Carries EAN-13 barcodes
- 12 digits + check digit
- The 978 and 979 prefixed are ISBN for books and publications
- First digits is GS1 Prefix

GTIN-14
- Not barcodes
- First digit indicates whether it's for packaging levels (1-8) or measures (9)
- GS1 Prefix starts at second digit
*/
package gtin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const GTIN_LENGTH = 14

type GTIN struct {
	Type   string
	Digits [GTIN_LENGTH]uint8
}

// The different GTIN types
const (
	GTIN8  string = "GTIN-8"  // 8 digits
	GTIN12 string = "GTIN-12" // 12 digits
	GTIN13 string = "GTIN-13" // 13 digits
	GTIN14 string = "GTIN-14" // 14 digits
)

const (
	EAN13   string = "EAN-13"
	EAN8    string = "EAN-8"
	UPCA    string = "UPC-A"
	ITF14   string = "ITF-14"
	UNKNOWN string = "UNKNOWN"
)

// String returns GTIN-14 as a string
func (gt GTIN) String() string {
	var s strings.Builder
	for _, m := range gt.Digits {
		s.WriteString(strconv.Itoa(int(m)))
	}
	return s.String()
}

// checkCheckDigit returns an error if the checkdigit is not valid
// https://www.gs1.org/services/how-calculate-check-digit-manually
// https://www.gs1us.org/tools/check-digit-calculator
func checkCheckDigit(gt GTIN) error {
	var multpliers = [GTIN_LENGTH - 1]uint8{3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3}
	var checksum int
	for n, m := range multpliers {
		checksum += int(gt.Digits[n] * m)
	}
	var checkdigit uint8
	if (checksum % 10) == 0 {
		// checksum is equal to a multiple of ten
		checkdigit = 0
	} else {
		//subtract from the higher multiple of ten
		checkdigit = uint8(int((checksum+10)/10)*10 - checksum)
	}

	if checkdigit != gt.Digits[GTIN_LENGTH-1] {
		return fmt.Errorf("invalid check digit")
	}

	return nil
}

// isRestrictedPrefix returns true if the GS1 prefix is restricted or a coupon code
func checkGS1Prefix(gt GTIN) error {

	if gt.Type != GTIN14 && gt.Type != GTIN13 {
		// Other GTIN types don't have GS1 prefixes
		return nil
	}

	var prefix int
	if gt.Type == GTIN14 {
		// In GTIN14, first digit is indicator
		prefix = 1
	}

	if gt.Digits[prefix] == 2 || (gt.Digits[prefix] == 0 && (gt.Digits[prefix+1] >= 2 && gt.Digits[prefix+1] <= 4)) {
		// Restricted prefixes 02, 04, or 2
		return errors.New("GS1 restricted prefix 02, 04 or 2")
	}
	if gt.Digits[prefix] == 9 && (gt.Digits[prefix+1] == 8 || gt.Digits[prefix+1] == 9) {
		// Coupon prefixes 98-99
		return errors.New("GS1 coupon prefix 98-99")
	}
	if gt.Digits[prefix] == 0 && gt.Digits[prefix+1] == 5 {
		// Coupon prefixes 05
		return errors.New("GS1 coupon prefix 05")
	}
	return nil
}

func (gt GTIN) Valid() bool {
	return checkCheckDigit(gt) == nil
}

func (gt GTIN) Legal() bool {
	return checkGS1Prefix(gt) == nil
}

// Carrier returns data carrier of the GTIN
func (gt GTIN) Carrier() string {

	var zeroes int
	for _, c := range gt.Digits {
		if c == 0 {
			zeroes++
		} else {
			break
		}
	}
	switch zeroes {
	case 0:
		return ITF14
	case 1:
		return EAN13
	case 2:
		return UPCA
	case 3:
		return UPCA
	case 4:
		return UPCA
	case 6:
		return EAN8
	}
	return UNKNOWN
}

// getGTINType returns the GTIN type based on length
func getGTINType(input string) (string, error) {
	switch len(input) {
	case 8:
		return GTIN8, nil
	case 12:
		return GTIN12, nil
	case 13:
		return GTIN13, nil
	case 14:
		return GTIN14, nil
	default:
		return "", fmt.Errorf("invalid length")
	}
}

// Atog converts a string to GTIN-14
// 1. Checks the GS1 prefix
// 2. Converts to 14 digits
// 3. Checks check digit
func Atog(input string) (GTIN, error) {

	var (
		gtin GTIN
		ch   byte
		pos  int
		curr int
	)
	var err error

	// Type
	gtin.Type, err = getGTINType(input)
	if err != nil {
		return gtin, err
	}

	curr = GTIN_LENGTH - len(input)
	for {
		if pos >= len(input) || curr >= GTIN_LENGTH {
			break
		}

		ch = input[pos]
		if '0' <= ch && ch <= '9' {
			gtin.Digits[curr] = ch - '0'
		} else if ch == 'X' {
			// Special case for ISBN
			gtin.Digits[curr] = 10
		} else {
			// we only accept numbers
			return gtin, fmt.Errorf("invalid digit")
		}

		pos++
		curr++
	}

	return gtin, nil
}
