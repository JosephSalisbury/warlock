package regression

import (
	"math"
	"testing"
)

func float64Equals(a, b float64) bool {
	return math.Abs(a-b) <= 0.001
}

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
				Start: 0,
				End:   1,

				Intercept: 0,
				Gradient:  1,

				Width: 0,
			},
		},
		"two negative samples": {
			input: []Sample{
				{-1, -1},
				{0, 0},
			},
			regression: &Regression{
				Start: -1,
				End:   0,

				Intercept: 0,
				Gradient:  1,

				Width: 0,
			},
		},
		"lower gradient line": {
			input: []Sample{
				{1, 1},
				{3, 2},
				{5, 3},
			},
			regression: &Regression{
				Start: 1,
				End:   5,

				Intercept: 0.5,
				Gradient:  0.5,

				Width: 0,
			},
		},
		"non-zero width line": {
			input: []Sample{
				{1, 2},
				{2, 1},
				{2, 3},
				{3, 2},
				{3, 4},
				{4, 3},
			},
			regression: &Regression{
				Start: 1,
				End:   4,

				Intercept: 1.363,
				Gradient:  0.454,

				Width: 1.142,
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

			if tc.regression != nil {
				if tc.regression.Start != regression.Start {
					t.Fatalf("expected Start: %v, got: %v", tc.regression.Start, regression.Start)
				}
				if tc.regression.End != regression.End {
					t.Fatalf("expected End: %v, got: %v", tc.regression.End, regression.End)
				}

				if !float64Equals(tc.regression.Intercept, regression.Intercept) {
					t.Fatalf("expected Intercept: %v, got: %v", tc.regression.Intercept, regression.Intercept)
				}
				if !float64Equals(tc.regression.Gradient, regression.Gradient) {
					t.Fatalf("expected Gradient: %v, got: %v", tc.regression.Gradient, regression.Gradient)
				}

				if !float64Equals(tc.regression.Width, regression.Width) {
					t.Fatalf("expected Width: %v, got: %v", tc.regression.Width, regression.Width)
				}
			}
		})
	}
}
