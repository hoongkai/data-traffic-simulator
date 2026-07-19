package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hoongkai/data-traffic-sim/internal/generator"
	"github.com/hoongkai/data-traffic-sim/internal/schema"
	"github.com/hoongkai/data-traffic-sim/internal/simulator"
)

func main() {
	dataset := os.Getenv("DATASET")

	if dataset == "" {
		fmt.Println("DATASET is required")
		return
	}

	dataSchema, err := schema.AnalyzeCSV(dataset)
	if err != nil {
		panic(err)
	}

	profile, err := schema.ProfileCSV(
		dataset,
		dataSchema,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Dataset: %s\n",
		dataset,
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

	fmt.Println()

	g := generator.New(
		time.Now().UnixNano(),
	)

	sim := simulator.New(
		g,
		profile,
		10,
	)

	stream := sim.Run(3 * time.Second)

	for row := range stream {
		fmt.Println(row)
	}
}
