package generator

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hoongkai/data-traffic-sim/internal/schema"
)

type Generator struct {
	rng *rand.Rand
}

func New(seed int64) *Generator {
	return &Generator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (g *Generator) Generate(profile schema.DatasetProfile) []any {
	row := make([]any, len(profile.Columns))

	for i, column := range profile.Columns {
		row[i] = g.generateColumn(column)
	}

	return row
}

func (g *Generator) generateColumn(column schema.ColumnProfile) any {
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

func (g *Generator) generateInteger(
	column schema.ColumnProfile,
) int64 {
	min := int64(math.Round(column.Min))
	max := int64(math.Round(column.Max))

	if min == max {
		return min
	}

	return min + g.rng.Int63n(max-min+1)
}

func (g *Generator) generateFloat(
	column schema.ColumnProfile,
) float64 {
	if column.Min == column.Max {
		return column.Min
	}

	return column.Min +
		g.rng.Float64()*
			(column.Max-column.Min)
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
