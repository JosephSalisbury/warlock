package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"

	"github.com/JosephSalisbury/warlock/regression"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

const (
	numRegressions = 10

	widthPadding = 25
)

func generateSineWaveSamples() []regression.Sample {
	s := []regression.Sample{}

	for i := 0; i < 360; i++ {
		x := float64(i)
		y := math.Sin((x * math.Pi) / 180)

		s = append(s, regression.Sample{X: x, Y: y})
	}

	return s
}

func samplesToContinuousSeries(samples []regression.Sample) chart.ContinuousSeries {
	xs := []float64{}
	ys := []float64{}

	for i := 0; i < len(samples); i++ {
		xs = append(xs, samples[i].X)
		ys = append(ys, samples[i].Y)
	}

	return chart.ContinuousSeries{
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorBlue,
		},

		XValues: xs,
		YValues: ys,
	}
}

func cutSamples(samples []regression.Sample) [][]regression.Sample {
	cutSamples := [][]regression.Sample{}
	sampleSize := len(samples) / numRegressions

	for i := 0; i < numRegressions; i++ {
		s := []regression.Sample{}

		for j := 0; j < sampleSize; j++ {
			s = append(s, samples[i*sampleSize+j])
		}

		cutSamples = append(cutSamples, s)
	}

	return cutSamples
}

func samplesToRegression(samples []regression.Sample) regression.Regression {
	regressionBuffer := regression.RegressionBuffer{}

	for _, sample := range samples {
		regressionBuffer.Add(sample)
	}

	regression, err := regressionBuffer.Regression()
	if err != nil {
		panic(err)
	}

	return *regression
}

func regressionToContinuousSeries(r regression.Regression) chart.ContinuousSeries {
	points := 10

	distance := r.End - r.Start
	step := distance / float64(points)

	xs := []float64{}
	ys := []float64{}

	for i := 0; i < points; i++ {
		x := r.Start + step*float64(i)
		y := r.Gradient*x + r.Intercept

		xs = append(xs, x)
		ys = append(ys, y)
	}

	return chart.ContinuousSeries{
		Style: chart.Style{
			Show:        true,
			StrokeWidth: r.Width * 2 * widthPadding,
			StrokeColor: drawing.ColorRed,
		},

		XValues: xs,
		YValues: ys,
	}
}

func main() {
	sineWaveSamples := generateSineWaveSamples()

	cutSineWaveSamples := cutSamples(sineWaveSamples)

	regressions := []regression.Regression{}
	for _, samples := range cutSineWaveSamples {
		regressions = append(regressions, samplesToRegression(samples))
	}

	regressionContinuousSerieses := []chart.Series{}
	for _, regression := range regressions {
		regressionContinuousSerieses = append(regressionContinuousSerieses, regressionToContinuousSeries(regression))
	}

	graph := chart.Chart{
		Series: []chart.Series{samplesToContinuousSeries(sineWaveSamples)},
	}
	graph.Series = append(graph.Series, regressionContinuousSerieses...)

	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		log.Fatalf(err.Error())
	}

	if err := ioutil.WriteFile("/Users/joe/Desktop/warlock-output.png", buffer.Bytes(), 0644); err != nil {
		log.Fatalf(err.Error())
	}
}
