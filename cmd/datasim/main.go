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
	dataSchema := schema.Schema{
		Columns: []schema.Column{
			{
				Name: "user_id",
				Type: schema.Integer,
			},
			{
				Name: "country",
				Type: schema.String,
			},
			{
				Name: "amount",
				Type: schema.Float,
			},
			{
				Name: "created_at",
				Type: schema.Timestamp,
			},
		},
	}

	if dataset := os.Getenv("DATASET"); dataset != "" {
		inferredSchema, err := schema.AnalyzeCSV(dataset)
		if err != nil {
			panic(err)
		}

		dataSchema = inferredSchema

		fmt.Println("Inferred schema:")

		for _, column := range dataSchema.Columns {
			fmt.Printf(
				"  %s: type=%s nullable=%t\n",
				column.Name,
				column.Type,
				column.Nullable,
			)
		}

		profile, err := schema.ProfileCSV(dataset, dataSchema)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nDataset profile:\n")
		fmt.Printf("  Rows: %d\n\n", profile.RowCount)

		for _, column := range profile.Columns {
			fmt.Printf("%s\n", column.Name)
			fmt.Printf("  Type: %s\n", column.Type)
			fmt.Printf("  Values: %d\n", column.Count)
			fmt.Printf("  Nulls: %d\n", column.NullCount)
			fmt.Printf("  Unique: %d\n", column.UniqueCount)

			if column.Type == schema.Integer ||
				column.Type == schema.Float {
				fmt.Printf("  Min: %.4f\n", column.Min)
				fmt.Printf("  Max: %.4f\n", column.Max)
				fmt.Printf("  Mean: %.4f\n", column.Mean)
			}

			if column.Type == schema.String {
				fmt.Println("  Frequencies:")

				for value, count := range column.Frequencies {
					fmt.Printf(
						"    %s: %d\n",
						value,
						count,
					)
				}
			}

			fmt.Println()
		}
	}

	g := generator.New(time.Now().UnixNano())

	sim := simulator.New(
		g,
		dataSchema,
		10,
	)

	stream := sim.Run(3 * time.Second)

	for row := range stream {
		fmt.Println(row)
	}
}
