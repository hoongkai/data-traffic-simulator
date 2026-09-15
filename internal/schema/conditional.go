package schema

import (
	"strconv"
	"strings"
)

func buildConditionalProfiles(
	records [][]string,
	profile DatasetProfile,
) []ConditionalProfile {
	var profiles []ConditionalProfile

	for conditionIndex, conditionColumn := range profile.Columns {

		if conditionColumn.Type != String {
			continue
		}

		conditionValues :=
			make(map[string]bool)

		for _, record := range records {
			if conditionIndex >= len(record) {
				continue
			}

			value := strings.TrimSpace(
				record[conditionIndex],
			)

			if value != "" {
				conditionValues[value] = true
			}
		}

		// A column with only one value does not provide
		// useful conditional information.
		if len(conditionValues) < 2 {
			continue
		}

		// Do not treat a column where every row has a
		// different value as a categorical dimension.
		if len(conditionValues) >= len(records) {
			continue
		}

		for conditionValue := range conditionValues {

			numeric :=
				make(map[string]NumericStats)

			for columnIndex, column := range profile.Columns {

				if column.Type != Integer &&
					column.Type != Float {
					continue
				}

				// IDs are excluded from conditional
				// statistical modelling.
				if column.IsID {
					continue
				}

				values := make([]float64, 0)

				for _, record := range records {

					if conditionIndex >= len(record) ||
						columnIndex >= len(record) {
						continue
					}

					recordCondition :=
						strings.TrimSpace(
							record[conditionIndex],
						)

					if recordCondition !=
						conditionValue {
						continue
					}

					value :=
						strings.TrimSpace(
							record[columnIndex],
						)

					number, err :=
						strconv.ParseFloat(
							value,
							64,
						)

					if err != nil {
						continue
					}

					values = append(
						values,
						number,
					)
				}

				if len(values) == 0 {
					continue
				}

				mean :=
					calculateMean(values)

				numeric[column.Name] =
					NumericStats{
						Count: len(values),
						Min: calculateMin(
							values,
						),
						Max: calculateMax(
							values,
						),
						Mean: mean,
						StdDev: calculateStdDev(
							values,
							mean,
						),
					}
			}

			if len(numeric) == 0 {
				continue
			}

			profiles = append(
				profiles,
				ConditionalProfile{
					ConditionColumn: conditionColumn.Name,
					ConditionValue:  conditionValue,
					Numeric:         numeric,
				},
			)
		}
	}

	return profiles
}

func calculateMin(
	values []float64,
) float64 {
	if len(values) == 0 {
		return 0
	}

	result := values[0]

	for _, value := range values[1:] {
		if value < result {
			result = value
		}
	}

	return result
}

func calculateMax(
	values []float64,
) float64 {
	if len(values) == 0 {
		return 0
	}

	result := values[0]

	for _, value := range values[1:] {
		if value > result {
			result = value
		}
	}

	return result
}
