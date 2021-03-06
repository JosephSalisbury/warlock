package regression

import (
	"errors"
	"math"
)

type Sample struct {
	X float64
	Y float64
}

type Regression struct {
	Start float64
	End   float64

	Intercept float64
	Gradient  float64

	Width float64
}

type RegressionBuffer struct {
	start float64
	end   float64

	n float64

	sx  float64
	sy  float64
	sxy float64
	sx2 float64
	sy2 float64

	width float64
}

func (r *RegressionBuffer) Add(s Sample) {
	if r.n == 0 {
		r.start = s.X
	}
	r.end = s.X

	r.n++

	r.sx = r.sx + s.X
	r.sy = r.sy + s.Y
	r.sxy = r.sxy + s.X*s.Y
	r.sx2 = r.sx2 + s.X*s.X
	r.sy2 = r.sy2 + s.Y*s.Y

	if r.n >= 2 {
		gradient, err := r.Gradient()
		if err != nil {
			panic(err)
		}
		intercept, err := r.Intercept()
		if err != nil {
			panic(err)
		}

		d := math.Abs(s.Y - (gradient*s.X + intercept))
		r.width = math.Max(r.width, d)
	}
}

func (r *RegressionBuffer) Intercept() (float64, error) {
	if r.n < 2 {
		return 0, errors.New("intercept requires at least two samples")
	}

	intercept := (r.sy*r.sx2 - r.sx*r.sxy) / (r.n*r.sx2 - r.sx*r.sx)

	return intercept, nil
}

func (r *RegressionBuffer) Gradient() (float64, error) {
	if r.n < 2 {
		return 0, errors.New("gradient requires at least two samples")
	}

	gradient := (r.n*r.sxy - r.sx*r.sy) / (r.n*r.sx2 - r.sx*r.sx)

	return gradient, nil
}

func (r *RegressionBuffer) Regression() (*Regression, error) {
	if r.n < 2 {
		return nil, errors.New("regression requires at least two samples")
	}

	intercept, err := r.Intercept()
	if err != nil {
		return nil, err
	}
	gradient, err := r.Gradient()
	if err != nil {
		return nil, err
	}

	regression := Regression{
		Start: r.start,
		End:   r.end,

		Intercept: intercept,
		Gradient:  gradient,

		Width: r.width,
	}

	return &regression, nil
}
