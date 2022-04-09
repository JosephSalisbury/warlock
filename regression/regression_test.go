package regression

import (
	"reflect"
	"testing"
)

func TestRegressionBuffer(t *testing.T) {
	tests := map[string]struct {
		input      []Sample
		regression *Regression
		err        string
	}{
		"zero samples": {
			input: []Sample{},
			err:   "regression requires at least two samples",
		},
		"one sample": {
			input: []Sample{
				{0, 0},
			},
			err: "regression requires at least two samples",
		},
		"two samples": {
			input: []Sample{
				{0, 0},
				{1, 1},
			},
			regression: &Regression{
				start: 0,
				end:   1,

				intercept: 0,
				gradient:  1,
			},
		},
		"two negative samples": {
			input: []Sample{
				{-1, -1},
				{0, 0},
			},
			regression: &Regression{
				start: -1,
				end:   0,

				intercept: 0,
				gradient:  1,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			regressionBuffer := RegressionBuffer{}

			for _, s := range tc.input {
				regressionBuffer.Add(s)
			}

			regression, err := regressionBuffer.Regression()

			if tc.err != "" && err == nil {
				t.Fatalf("expected error: %v, got no error", tc.err)
			}

			if tc.err != "" && tc.err != err.Error() {
				t.Fatalf("expected error: %v, got: %v", tc.err, err.Error())
			}

			if !reflect.DeepEqual(tc.regression, regression) {
				t.Fatalf("expected: %v, got: %v", tc.regression, regression)
			}
		})
	}
}
