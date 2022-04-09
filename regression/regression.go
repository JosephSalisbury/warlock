package regression

import "errors"

type Sample struct {
	x float64
	y float64
}

type Regression struct {
	Start float64
	End   float64

	Intercept float64
	Gradient  float64
}

type RegressionBuffer struct {
	n float64

	start float64
	end   float64

	sx  float64
	sy  float64
	sxy float64
	sx2 float64
	sy2 float64
}

func (r *RegressionBuffer) Add(s Sample) {
	if r.n == 0 {
		r.start = s.x
	}
	r.end = s.x

	r.n++

	r.sx = r.sx + s.x
	r.sy = r.sy + s.y
	r.sxy = r.sxy + s.x*s.y
	r.sx2 = r.sx2 + s.x*s.x
	r.sy2 = r.sy2 + s.y*s.y
}

func (r *RegressionBuffer) Regression() (*Regression, error) {
	if r.n < 2 {
		return nil, errors.New("regression requires at least two samples")
	}

	intercept := (r.sy*r.sx2 - r.sx*r.sxy) / (r.n*r.sx2 - r.sx*r.sx)
	gradient := (r.n*r.sxy - r.sx*r.sy) / (r.n*r.sx2 - r.sx*r.sx)

	regression := Regression{
		Start: r.start,
		End:   r.end,

		Intercept: intercept,
		Gradient:  gradient,
	}

	return &regression, nil
}
