package simulator

import (
	"time"

	"github.com/hoongkai/data-traffic-sim/internal/generator"
	"github.com/hoongkai/data-traffic-sim/internal/schema"
)

type Simulator struct {
	generator *generator.Generator
	profile   schema.DatasetProfile
	rate      int
}

func New(
	g *generator.Generator,
	profile schema.DatasetProfile,
	rate int,
) *Simulator {
	return &Simulator{
		generator: g,
		profile:   profile,
		rate:      rate,
	}
}

func (s *Simulator) Run(
	duration time.Duration,
) <-chan []any {
	stream := make(chan []any)

	go func() {
		defer close(stream)

		ticker := time.NewTicker(
			time.Second / time.Duration(s.rate),
		)
		defer ticker.Stop()

		timeout := time.After(duration)

		for {
			select {
			case <-ticker.C:
				stream <- s.generator.Generate(s.profile)

			case <-timeout:
				return
			}
		}
	}()

	return stream
}
