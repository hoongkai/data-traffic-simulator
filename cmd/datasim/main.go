package main

import (
	"fmt"
	"os"

	"github.com/hoongkai/data-traffic-sim/internal/generator"
	"github.com/hoongkai/data-traffic-sim/internal/schema"
)

func main() {
	datasetPath := os.Getenv("DATASET")

	if datasetPath == "" {
		datasetPath = "testdata/Iris.csv"
	}

	dataSchema := schema.Schema{
		Columns: []schema.Column{
			{Name: "Id", Type: schema.Integer},
			{Name: "SepalLengthCm", Type: schema.Float},
			{Name: "SepalWidthCm", Type: schema.Float},
			{Name: "PetalLengthCm", Type: schema.Float},
			{Name: "PetalWidthCm", Type: schema.Float},
			{Name: "Species", Type: schema.String},
		},
	}

	profile, err := schema.ProfileCSV(
		datasetPath,
		dataSchema,
	)

	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Dataset: %s\n",
		datasetPath,
	)

	fmt.Printf(
		"Rows: %d\n\n",
		profile.RowCount,
	)

	// Display column profiles.
	for _, column := range profile.Columns {
		fmt.Printf(
			"%s: type=%s values=%d nulls=%d unique=%d",
			column.Name,
			column.Type,
			column.Count,
			column.NullCount,
			column.UniqueCount,
		)

		if column.IsID {
			fmt.Print(" [ID]")
		}

		fmt.Println()

		if column.Type == schema.Integer ||
			column.Type == schema.Float {

			fmt.Printf(
				"  min=%.4f max=%.4f mean=%.4f stddev=%.4f\n",
				column.Min,
				column.Max,
				column.Mean,
				column.StdDev,
			)
		}

		if column.Type == schema.String {
			fmt.Println("  frequencies:")

			for value, count := range column.Frequencies {
				fmt.Printf(
					"    %s: %d\n",
					value,
					count,
				)
			}
		}
	}

	// Display correlations between numeric columns.
	// ID columns are excluded.
	correlations := schema.CalculateCorrelations(
		profile.NumericData,
	)

	fmt.Println("\nCorrelations:")

	for _, column := range profile.Columns {
		if (column.Type == schema.Integer ||
			column.Type == schema.Float) &&
			!column.IsID {

			fmt.Printf(
				"%-18s",
				column.Name,
			)
		}
	}

	fmt.Println()

	for i := range correlations {
		for j := range correlations[i] {
			fmt.Printf(
				"%-18.4f",
				correlations[i][j],
			)
		}

		fmt.Println()
	}

	// Display conditional profiles.
	// These describe numeric distributions under
	// different categorical conditions.
	fmt.Println("\nConditional Profiles:")

	for _, conditional := range profile.ConditionalProfiles {

		fmt.Printf(
			"\n%s = %s\n",
			conditional.ConditionColumn,
			conditional.ConditionValue,
		)

		for columnName, stats := range conditional.Numeric {

			fmt.Printf(
				"  %-18s min=%.4f max=%.4f mean=%.4f stddev=%.4f\n",
				columnName,
				stats.Min,
				stats.Max,
				stats.Mean,
				stats.StdDev,
			)
		}
	}

	// Generate sample records, global numeric config
	fmt.Println("\nGenerated Records:")

	g := generator.New(42)

	for i := 0; i < 30; i++ {
		row := g.Generate(profile)
		fmt.Println(row)
	}
}
