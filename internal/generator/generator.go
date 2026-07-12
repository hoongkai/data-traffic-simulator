package generator

import (
	"fmt"
	"math/rand"
	"strconv"
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

func (g *Generator) GenerateValue(column schema.Column) string {
	switch column.Type {
	case schema.Integer:
		return strconv.Itoa(g.rng.Intn(1000000))

	case schema.Float:
		return fmt.Sprintf("%.2f", g.rng.Float64()*1000)

	case schema.Boolean:
		return strconv.FormatBool(g.rng.Intn(2) == 1)

	case schema.Timestamp:
		return time.Now().UTC().Format(time.RFC3339Nano)

	case schema.String:
		return fmt.Sprintf("value_%d", g.rng.Intn(1000000))

	default:
		return ""
	}
}

func (g *Generator) GenerateRow(s schema.Schema) []string {
	row := make([]string, len(s.Columns))

	for i, column := range s.Columns {
		row[i] = g.GenerateValue(column)
	}

	return row
}
