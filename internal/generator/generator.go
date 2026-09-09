package generator

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hoongkai/data-traffic-sim/internal/schema"
)

type Generator struct {
	rng    *rand.Rand
	nextID map[string]int64
}

func New(seed int64) *Generator {
	return &Generator{
		rng:    rand.New(rand.NewSource(seed)),
		nextID: make(map[string]int64),
	}
}

func (g *Generator) Generate(
	profile schema.DatasetProfile,
) []any {
	row := make([]any, len(profile.Columns))

	for i, column := range profile.Columns {
		row[i] = g.generateColumn(column)
	}

	return row
}

func (g *Generator) generateColumn(
	column schema.ColumnProfile,
) any {
	if column.Count == 0 {
		return nil
	}

	if column.NullCount > 0 {
		nullProbability :=
			float64(column.NullCount) /
				float64(column.Count+column.NullCount)

		if g.rng.Float64() < nullProbability {
			return nil
		}
	}

	switch column.Type {
	case schema.Integer:
		if column.IsID {
			return g.generateID(column.Name)
		}

		return g.generateInteger(column)

	case schema.Float:
		return g.generateFloat(column)

	case schema.Boolean:
		return g.rng.Intn(2) == 1

	case schema.Timestamp:
		return g.generateTimestamp()

	case schema.String:
		return g.generateCategorical(column)

	default:
		return nil
	}
}

func (g *Generator) generateID(
	columnName string,
) int64 {
	if _, exists := g.nextID[columnName]; !exists {
		g.nextID[columnName] = 1
	}

	id := g.nextID[columnName]
	g.nextID[columnName]++

	return id
}

func (g *Generator) generateInteger(
	column schema.ColumnProfile,
) int64 {
	value := g.generateGaussian(
		column.Mean,
		column.StdDev,
	)

	value = clamp(
		value,
		column.Min,
		column.Max,
	)

	return int64(math.Round(value))
}

func (g *Generator) generateFloat(
	column schema.ColumnProfile,
) float64 {
	value := g.generateGaussian(
		column.Mean,
		column.StdDev,
	)

	return clamp(
		value,
		column.Min,
		column.Max,
	)
}

func (g *Generator) generateGaussian(
	mean float64,
	stdDev float64,
) float64 {
	if stdDev <= 0 {
		return mean
	}

	/*
		Box-Muller transform.

		Converts two uniformly distributed random
		numbers into a normally distributed value.
	*/

	u1 := g.rng.Float64()

	// Avoid ln(0).
	if u1 == 0 {
		u1 = math.SmallestNonzeroFloat64
	}

	u2 := g.rng.Float64()

	z := math.Sqrt(
		-2.0 * math.Log(u1),
	) * math.Cos(
		2.0 * math.Pi * u2,
	)

	return mean + z*stdDev
}

func clamp(
	value float64,
	min float64,
	max float64,
) float64 {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

func (g *Generator) generateCategorical(
	column schema.ColumnProfile,
) string {
	if len(column.Frequencies) == 0 {
		return fmt.Sprintf(
			"value_%d",
			g.rng.Int63(),
		)
	}

	total := 0

	for _, count := range column.Frequencies {
		total += count
	}

	target := g.rng.Intn(total)

	current := 0

	for value, count := range column.Frequencies {
		current += count

		if target < current {
			return value
		}
	}

	return ""
}

func (g *Generator) generateTimestamp() time.Time {
	now := time.Now()

	offset := time.Duration(
		g.rng.Int63n(int64(24 * time.Hour)),
	)

	return now.Add(-offset)
}