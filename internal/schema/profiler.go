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

		columnProfile.IsID = isIdentifier(columnProfile)

		profile.Columns[i] = columnProfile
	}

	return profile, nil
}

func isIdentifier(column ColumnProfile) bool {
	if column.Type != Integer {
		return false
	}

	name := strings.ToLower(strings.TrimSpace(column.Name))

	if name == "id" ||
		strings.HasSuffix(name, "_id") ||
		strings.HasSuffix(name, "id") {
		return column.UniqueCount == column.Count &&
			column.NullCount == 0
	}

	return false
}
