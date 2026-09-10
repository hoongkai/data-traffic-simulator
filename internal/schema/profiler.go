package schema

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

func ProfileCSV(
	path string,
	dataSchema Schema,
) (DatasetProfile, error) {
	file, err := os.Open(path)
	if err != nil {
		return DatasetProfile{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return DatasetProfile{}, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return DatasetProfile{}, err
	}

	profile := DatasetProfile{
		RowCount: len(records),
		Columns:  make([]ColumnProfile, len(dataSchema.Columns)),
	}

	for i, column := range dataSchema.Columns {
		columnProfile := ColumnProfile{
			Name:        column.Name,
			Type:        column.Type,
			Frequencies: make(map[string]int),
		}

		values := make([]string, 0, len(records))
		numericValues := make([]float64, 0, len(records))

		for _, record := range records {
			if i >= len(record) {
				columnProfile.NullCount++
				continue
			}

			value := strings.TrimSpace(record[i])

			if value == "" {
				columnProfile.NullCount++
				continue
			}

			columnProfile.Count++
			columnProfile.Frequencies[value]++
			values = append(values, value)

			switch column.Type {
			case Integer, Float:
				number, err := strconv.ParseFloat(value, 64)
				if err != nil {
					continue
				}

				numericValues = append(
					numericValues,
					number,
				)

				if len(numericValues) == 1 ||
					number < columnProfile.Min {
					columnProfile.Min = number
				}

				if len(numericValues) == 1 ||
					number > columnProfile.Max {
					columnProfile.Max = number
				}
			}
		}

		columnProfile.UniqueCount =
			len(columnProfile.Frequencies)

		if len(numericValues) > 0 {
			columnProfile.Mean =
				calculateMean(numericValues)

			columnProfile.StdDev =
				calculateStdDev(
					numericValues,
					columnProfile.Mean,
				)
		}

		columnProfile.IsID = isIdentifier(
			columnProfile,
			values,
		)

		profile.Columns[i] = columnProfile
	}

	// Collect numeric data for correlation analysis.
	numericColumns := make([]int, 0)

	for i, column := range dataSchema.Columns {
		if column.Type == Integer ||
			column.Type == Float {
			numericColumns = append(
				numericColumns,
				i,
			)
		}
	}

	for _, record := range records {
		row := make(
			[]float64,
			0,
			len(numericColumns),
		)

		for _, columnIndex := range numericColumns {
			if columnIndex >= len(record) {
				continue
			}

			value := strings.TrimSpace(
				record[columnIndex],
			)

			number, err := strconv.ParseFloat(
				value,
				64,
			)

			if err != nil {
				continue
			}

			row = append(row, number)
		}

		if len(row) == len(numericColumns) {
			profile.NumericData =
				append(profile.NumericData, row)
		}
	}

	return profile, nil
}

func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum float64

	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}

func calculateStdDev(
	values []float64,
	mean float64,
) float64 {
	if len(values) < 2 {
		return 0
	}

	var sumSquaredDiff float64

	for _, value := range values {
		diff := value - mean
		sumSquaredDiff += diff * diff
	}

	variance :=
		sumSquaredDiff /
			float64(len(values))

	return sqrt(variance)
}

func sqrt(value float64) float64 {
	if value <= 0 {
		return 0
	}

	x := value

	for i := 0; i < 10; i++ {
		x = 0.5 * (x + value/x)
	}

	return x
}

func isIdentifier(
	column ColumnProfile,
	values []string,
) bool {
	if column.NullCount > 0 || column.Count == 0 {
		return false
	}

	uniqueness :=
		float64(column.UniqueCount) /
			float64(column.Count)

	name :=
		strings.ToLower(
			strings.TrimSpace(column.Name),
		)

	if name == "id" ||
		strings.HasSuffix(name, "_id") ||
		strings.HasSuffix(name, "id") {
		return uniqueness >= 0.95
	}

	if strings.HasSuffix(name, "_number") ||
		strings.HasSuffix(name, "_no") ||
		strings.HasSuffix(name, "_code") {
		return uniqueness >= 0.95
	}

	if strings.Contains(name, "number") ||
		strings.Contains(name, "identifier") ||
		strings.Contains(name, "key") ||
		strings.Contains(name, "code") {
		return uniqueness >= 0.95 &&
			isSequential(values)
	}

	return false
}

func isSequential(values []string) bool {
	if len(values) < 2 {
		return false
	}

	previous, err := strconv.ParseInt(
		strings.TrimSpace(values[0]),
		10,
		64,
	)

	if err != nil {
		return false
	}

	var direction int64

	for _, value := range values[1:] {
		current, err := strconv.ParseInt(
			strings.TrimSpace(value),
			10,
			64,
		)

		if err != nil {
			return false
		}

		difference := current - previous

		if difference == 0 {
			return false
		}

		if direction == 0 {
			direction = difference
		}

		if difference != direction {
			return false
		}

		previous = current
	}

	return true
}
