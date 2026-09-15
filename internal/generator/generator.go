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

	categoricalValues := make(map[string]string)

	for i, column := range profile.Columns {
		if column.IsID {
			row[i] = g.generateID(column.Name)
			continue
		}

		if column.Type == schema.String {
			value := g.generateCategorical(column)

			row[i] = value
			categoricalValues[column.Name] = value

			continue
		}
	}

	for i, column := range profile.Columns {
		if column.IsID ||
			column.Type == schema.String {
			continue
		}

		row[i] = g.generateColumn(
			column,
			profile,
			categoricalValues,
		)
	}

	return row
}

func (g *Generator) generateColumn(
	column schema.ColumnProfile,
	profile schema.DatasetProfile,
	categoricalValues map[string]string,
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

	conditionalStats := findConditionalStats(
		column.Name,
		profile.ConditionalProfiles,
		categoricalValues,
	)

	if conditionalStats != nil {
		switch column.Type {
		case schema.Integer:
			return g.generateConditionalInteger(*conditionalStats)

		case schema.Float:
			return g.generateConditionalFloat(*conditionalStats)
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

	default:
		return nil
	}
}

func findConditionalStats(
	columnName string,
	profiles []schema.ConditionalProfile,
	categoricalValues map[string]string,
) *schema.NumericStats {
	for _, profile := range profiles {
		generatedValue, exists :=
			categoricalValues[profile.ConditionColumn]

		if !exists {
			continue
		}

		if generatedValue != profile.ConditionValue {
			continue
		}

		stats, exists := profile.Numeric[columnName]

		if !exists {
			continue
		}

		return &stats
	}

	return nil
}

func (g *Generator) generateConditionalInteger(
	stats schema.NumericStats,
) int64 {
	value := g.generateGaussian(
		stats.Mean,
		stats.StdDev,
	)

	value = clamp(
		value,
		stats.Min,
		stats.Max,
	)

	return int64(math.Round(value))
}

func (g *Generator) generateConditionalFloat(
	stats schema.NumericStats,
) float64 {
	value := g.generateGaussian(
		stats.Mean,
		stats.StdDev,
	)

	return clamp(
		value,
		stats.Min,
		stats.Max,
	)
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

	// below uses Box-Muller transform:
	// converts two uniformly distributed random numbers into a normally distributed value

	u1 := g.rng.Float64()

	// avoid ln(0)
	if u1 == 0 {
		u1 = math.SmallestNonzeroFloat64
	}

	u2 := g.rng.Float64()

	z := math.Sqrt(
		-2.0*math.Log(u1),
	) * math.Cos(
		2.0*math.Pi*u2,
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
