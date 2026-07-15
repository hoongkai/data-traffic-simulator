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

	// Added dataset reading for analyzer.go use
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

		fmt.Println()
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
