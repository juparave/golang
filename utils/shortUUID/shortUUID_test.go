package shortuuid

import "testing"

type Tests struct {
	input    int // wanted length
	expected int // expected length
}

func TestShortUUID(t *testing.T) {

	tests := []Tests{
		{3, 3},
		{4, 4},
		{5, 5},
		{6, 6},
		{7, 7},
		{8, 8},
		{9, 9},
		{10, 10},
		{11, 11},
		{12, 12},
		{13, 13},
		{14, 14},
	}

	for _, test := range tests {
		actual := ShortUUID(test.input)
		if len(actual) != test.expected {
			t.Errorf("ShortUUID(%d) = %d, expected %d", test.input, len(actual), test.expected)
		}
	}

}
