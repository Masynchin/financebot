package main

import (
	"bytes"
	"sort"

	"github.com/wcharczuk/go-chart/v2"
)

// RenderExpencesGroupedByCategory renders user expences with pie chart scheme
func RenderExpencesGroupedByCategory(expences []TotalExpence) (bytes.Buffer, error) {
	sortExpencesByAmount(expences)
	chartValues := expencesToChartValues(expences)
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: chartValues,
	}

	var b bytes.Buffer
	err := pie.Render(chart.PNG, &b)

	return b, err
}

func expencesToChartValues(expences []TotalExpence) []chart.Value {
	chartValues := make([]chart.Value, len(expences))
	for _, exp := range expences {
		chartValues = append(chartValues, chart.Value{Value: float64(exp.Amount), Label: exp.Category})
	}

	return chartValues
}

func sortExpencesByAmount(expences []TotalExpence) {
	sort.Slice(expences, func(i, j int) bool { return expences[i].Amount < expences[j].Amount })
}
