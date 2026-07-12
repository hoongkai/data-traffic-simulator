package simulator

import (
	"time"

	"github.com/hoongkai/data-traffic-sim/internal/generator"
	"github.com/hoongkai/data-traffic-sim/internal/schema"
)

type Simulator struct {
	generator *generator.Generator
	schema    schema.Schema
	rate      int
}

func New(
	g *generator.Generator,
	s schema.Schema,
	rate int,
) *Simulator {
	return &Simulator{
		generator: g,
		schema:    s,
		rate:      rate,
	}
}

func (s *Simulator) Run(
	duration time.Duration,
) <-chan []string {
	output := make(chan []string)

	go func() {
		defer close(output)

		ticker := time.NewTicker(time.Second / time.Duration(s.rate))
		defer ticker.Stop()

		timer := time.NewTimer(duration)
		defer timer.Stop()

		for {
			select {
			case <-ticker.C:
				output <- s.generator.GenerateRow(s.schema)

			case <-timer.C:
				return
			}
		}
	}()

	return output
}
