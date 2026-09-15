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

		if !isCategoricalCandidate(conditionColumn) {
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

		for conditionValue := range conditionValues {

			numeric :=
				make(map[string]NumericStats)

			numericColumns :=
				make([]int, 0)

			for columnIndex, column := range profile.Columns {

				if column.Type != Integer &&
					column.Type != Float {
					continue
				}

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

				numericColumns =
					append(
						numericColumns,
						columnIndex,
					)
			}

			if len(numeric) == 0 {
				continue
			}

			correlations :=
				buildConditionalCorrelations(
					records,
					profile,
					conditionIndex,
					conditionValue,
					numericColumns,
				)

			profiles = append(
				profiles,
				ConditionalProfile{
					ConditionColumn: conditionColumn.Name,
					ConditionValue:  conditionValue,
					Numeric:         numeric,
					Correlations:    correlations,
				},
			)
		}
	}

	return profiles
}

func buildConditionalCorrelations(
	records [][]string,
	profile DatasetProfile,
	conditionIndex int,
	conditionValue string,
	numericColumns []int,
) map[string]map[string]float64 {

	values := make([][]float64, 0)

	for _, record := range records {
		if conditionIndex >= len(record) {
			continue
		}

		recordCondition :=
			strings.TrimSpace(
				record[conditionIndex],
			)

		if recordCondition != conditionValue {
			continue
		}

		row := make(
			[]float64,
			0,
			len(numericColumns),
		)

		valid := true

		for _, columnIndex := range numericColumns {

			if columnIndex >= len(record) {
				valid = false
				break
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
				valid = false
				break
			}

			row = append(row, number)
		}

		if valid {
			values = append(values, row)
		}
	}

	matrix :=
		CalculateCorrelations(values)

	result :=
		make(map[string]map[string]float64)

	for i, columnI := range numericColumns {

		nameI :=
			profile.Columns[columnI].Name

		result[nameI] =
			make(map[string]float64)

		for j, columnJ := range numericColumns {

			nameJ :=
				profile.Columns[columnJ].Name

			if i < len(matrix) &&
				j < len(matrix[i]) {

				result[nameI][nameJ] =
					matrix[i][j]
			}
		}
	}

	return result
}

func isCategoricalCandidate(
	column ColumnProfile,
) bool {
	if column.Type != String {
		return false
	}

	if column.Count == 0 ||
		column.UniqueCount < 2 {
		return false
	}

	uniqueRatio :=
		float64(column.UniqueCount) /
			float64(column.Count)

	if uniqueRatio > 0.2 {
		return false
	}

	if column.UniqueCount > 100 {
		return false
	}

	return true
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
