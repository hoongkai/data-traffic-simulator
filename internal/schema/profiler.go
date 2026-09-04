package schema

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

func ProfileCSV(path string, dataSchema Schema) (DatasetProfile, error) {
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

		var sum float64
		var numericCount int

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

				if numericCount == 0 || number < columnProfile.Min {
					columnProfile.Min = number
				}

				if numericCount == 0 || number > columnProfile.Max {
					columnProfile.Max = number
				}

				sum += number
				numericCount++
			}
		}

		columnProfile.UniqueCount = len(columnProfile.Frequencies)

		if numericCount > 0 {
			columnProfile.Mean = sum / float64(numericCount)
		}

		columnProfile.IsID = isIdentifier(
			columnProfile,
			values,
		)

		profile.Columns[i] = columnProfile
	}

	return profile, nil
}

func isIdentifier(
	column ColumnProfile,
	values []string,
) bool {
	if column.NullCount > 0 || column.Count == 0 {
		return false
	}

	uniqueness := float64(column.UniqueCount) /
		float64(column.Count)

	name := strings.ToLower(
		strings.TrimSpace(column.Name),
	)

	// Identifiers
	if name == "id" ||
		strings.HasSuffix(name, "_id") ||
		strings.HasSuffix(name, "id") {
		return uniqueness >= 0.95
	}

	// Other identifiers column name variations
	if strings.HasSuffix(name, "_number") ||
		strings.HasSuffix(name, "_no") ||
		strings.HasSuffix(name, "_code") {
		return uniqueness >= 0.95
	}

	// Only do sequential when column name implies it is an identifier
	// TODO: perhaps improve on string/UUID detection
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
