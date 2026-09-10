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
			{
				Name: "Id",
				Type: schema.Integer,
			},
			{
				Name: "SepalLengthCm",
				Type: schema.Float,
			},
			{
				Name: "SepalWidthCm",
				Type: schema.Float,
			},
			{
				Name: "PetalLengthCm",
				Type: schema.Float,
			},
			{
				Name: "PetalWidthCm",
				Type: schema.Float,
			},
			{
				Name: "Species",
				Type: schema.String,
			},
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

	for _, column := range profile.Columns {
		fmt.Printf(
			"%s: type=%s values=%d nulls=%d unique=%d\n",
			column.Name,
			column.Type,
			column.Count,
			column.NullCount,
			column.UniqueCount,
		)

		if column.Type == schema.Integer ||
			column.Type == schema.Float {
			fmt.Printf(
				"  min=%.4f max=%.4f mean=%.4f\n",
				column.Min,
				column.Max,
				column.Mean,
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

	fmt.Println()

	g := generator.New(42)

	for i := 0; i < 30; i++ {
		row := g.Generate(profile)

		fmt.Println(row)
	}
}
