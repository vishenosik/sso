package services

import (
	"log"

	"github.com/vishenosik/sso/internal/entities"
)

type System struct {
	metrics entities.Metrics
}

func NewSystem(
	p1 int64,
	p2 bool,
	p3 int,
) *System {
	return &System{
		metrics: entities.Metrics{
			&entities.Metric{
				P1: p1,
				P2: p2,
				P3: p3,
			},
			&entities.Metric{
				P1: p1,
				P2: true,
				P3: p3,
			},
			&entities.Metric{
				P1: p1,
				P2: false,
			},
		},
	}
}

func (s *System) LogMetrics(metrics ...*entities.Metric) error {
	s.metrics = append(s.metrics, metrics...)
	log.Print("logging metrics: ", s.metrics)
	return nil
}

func (s *System) FetchMetrics() entities.Metrics {
	return s.metrics
}
