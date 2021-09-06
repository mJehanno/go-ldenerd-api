package goldmanager

import "testing"

type coinTest struct {
	val      int
	src      Currency
	dest     Currency
	expected int
}

var data = []coinTest{
	{10, 3, 0, 1000},
	{10, 3, 1, 100},
	{10, 3, 2, 20},
	{10, 3, 4, 1},
	{10, 0, 1, 1},
	{100, 0, 2, 2},
	{100, 0, 3, 1},
	{1000, 0, 4, 1},
	{1, 2, 0, 50},
}

func TestConvert(t *testing.T) {
	for _, d := range data {
		if r := Convert(d.val, d.src, d.dest); r != d.expected {
			t.Errorf("Conversion of %d %s to %s was incorrect, got: %d, wanted: %d.", d.val, d.src, d.dest, r, d.expected)
		}
	}
}
