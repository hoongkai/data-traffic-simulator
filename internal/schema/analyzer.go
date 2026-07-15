package schema

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
)

func AnalyzeCSV(path string) (Schema, error) {
	file, err := os.Open(path)
	if err != nil {
		return Schema{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return Schema{}, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return Schema{}, err
	}

	columns := make([]Column, len(header))

	for i, name := range header {
		values := make([]string, 0, len(records))
		nullable := false

		for _, record := range records {
			if i >= len(record) {
				nullable = true
				continue
			}

			value := strings.TrimSpace(record[i])

			if value == "" {
				nullable = true
				continue
			}

			values = append(values, value)
		}

		columns[i] = Column{
			Name:     name,
			Type:     inferType(values),
			Nullable: nullable,
		}
	}

	return Schema{
		Columns: columns,
	}, nil
}

func inferType(values []string) DataType {
	if len(values) == 0 {
		return String
	}

	allInteger := true
	allFloat := true
	allBoolean := true
	allTimestamp := true

	for _, value := range values {
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			allInteger = false
		}

		if _, err := strconv.ParseFloat(value, 64); err != nil {
			allFloat = false
		}

		if _, err := strconv.ParseBool(value); err != nil {
			allBoolean = false
		}

		if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
			allTimestamp = false
		}
	}

	switch {
	case allInteger:
		return Integer
	case allFloat:
		return Float
	case allBoolean:
		return Boolean
	case allTimestamp:
		return Timestamp
	default:
		return String
	}
}
