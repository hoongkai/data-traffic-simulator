package schema

import "math"

func CalculateCorrelations(
	values [][]float64,
) [][]float64 {
	if len(values) == 0 {
		return nil
	}

	columnCount := len(values[0])

	correlations := make(
		[][]float64,
		columnCount,
	)

	for i := range correlations {
		correlations[i] = make(
			[]float64,
			columnCount,
		)
	}

	for i := 0; i < columnCount; i++ {
		for j := 0; j < columnCount; j++ {
			correlations[i][j] =
				pearsonCorrelation(
					getColumn(values, i),
					getColumn(values, j),
				)
		}
	}

	return correlations
}

func getColumn(
	values [][]float64,
	column int,
) []float64 {
	result := make(
		[]float64,
		0,
		len(values),
	)

	for _, row := range values {
		if column < len(row) {
			result = append(result, row[column])
		}
	}

	return result
}

func pearsonCorrelation(
	x []float64,
	y []float64,
) float64 {
	if len(x) != len(y) || len(x) < 2 {
		return 0
	}

	var sumX float64
	var sumY float64

	for i := range x {
		sumX += x[i]
		sumY += y[i]
	}

	meanX := sumX / float64(len(x))
	meanY := sumY / float64(len(y))

	var numerator float64
	var sumSquaredX float64
	var sumSquaredY float64

	for i := range x {
		diffX := x[i] - meanX
		diffY := y[i] - meanY

		numerator += diffX * diffY
		sumSquaredX += diffX * diffX
		sumSquaredY += diffY * diffY
	}

	denominator := math.Sqrt(
		sumSquaredX * sumSquaredY,
	)

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}
