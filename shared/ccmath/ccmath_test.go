package ccmath

import (
	"testing"
)

func TestSolving(t *testing.T) {
	for i, test := range []struct {
		postFixInput string
		output       float64
	}{
		{"14 -5 /", -2.8},
		{"20 3 -4 + *", -20.0},
		{"2 3 + 5 / 4 5 - *", -1.0},
	} {
		result, err := Solve(test.postFixInput)
		if err != nil {
			t.Errorf(`Test %d: Solve("%s") returned error: %s`, i, test.postFixInput, err)
		}
		if result != test.output {
			t.Errorf(`Test %d: Solve("%s") returned value %f, was expecting %f`, i, test.postFixInput, result, test.output)
		}
	}
}
