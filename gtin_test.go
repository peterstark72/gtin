package gtin

import (
	"fmt"
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		got GTIN
	}{
		{GTIN{GTIN13, [GTIN_LENGTH]uint8{0, 9, 7, 8, 0, 6, 7, 0, 0, 2, 2, 1, 5, 1}}},
		{GTIN{GTIN13, [GTIN_LENGTH]uint8{0, 0, 0, 4, 1, 2, 5, 0, 5, 0, 0, 7, 3, 5}}},
		{GTIN{GTIN13, [GTIN_LENGTH]uint8{0, 9, 7, 8, 1, 9, 4, 9, 0, 0, 3, 7, 2, 7}}},
		{GTIN{GTIN13, [GTIN_LENGTH]uint8{0, 9, 7, 8, 0, 5, 9, 3, 2, 9, 7, 0, 6, 3}}},
		{GTIN{GTIN13, [GTIN_LENGTH]uint8{0, 9, 2, 9, 1, 0, 4, 1, 5, 0, 0, 2, 1, 0}}},
	}

	for _, tt := range tests {
		v := checkCheckDigit(tt.got)
		if v != nil {
			t.Errorf("wrong result")
		}
	}
}

func TestAtog(t *testing.T) {

	tests := []struct {
		got  string
		want string
	}{
		{"614141000012", "GTIN-12 00614141000012"},
		{"00614141000029", "GTIN-14 00614141000029"},
		{"614141000777", "GTIN-12 00614141000777"},
		{"50614141000994", "GTIN-14 50614141000994"},
	}

	for _, tt := range tests {
		result, err := Atog(tt.got)
		if err != nil {
			t.Error(err)
		}
		if tt.want != result.String() {
			t.Errorf("wanted %v, got %v", tt.want, result)
		}
	}
}

func TestGetCode(t *testing.T) {

	c, _ := Atog("08719076050360")
	fmt.Println(c.Carrier())

}
