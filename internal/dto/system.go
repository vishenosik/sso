package dto

import (
	"fmt"
	"strconv"

	"github.com/vishenosik/sso-sdk/api"
	"github.com/vishenosik/sso/internal/entities"
	"github.com/vishenosik/sso/internal/services"
	"golang.org/x/exp/constraints"
)

type SystemDTO struct {
	service *services.System
}

func NewSystemDTO(service *services.System) *SystemDTO {
	return &SystemDTO{
		service: service,
	}
}

func metricsEntitiesToApi(metrics entities.Metrics) api.Metrics {
	out := make(api.Metrics, 0, len(metrics))
	for _, m := range metrics {
		out = append(out, &api.Metric{
			Param1: fmt.Sprintf("%d", m.P1),
			Param2: fmt.Sprint(m.P2),
			Param3: fmt.Sprintf("%d", m.P3),
		})
	}
	return out
}

func metricsApiToEntities(metrics api.Metrics) entities.Metrics {
	out := make(entities.Metrics, 0, len(metrics))
	for _, m := range metrics {
		out = append(out, &entities.Metric{
			P1: atoi[int64](m.Param1),
			P2: parseBool(m.Param2),
			P3: atoi[int](m.Param3),
		})
	}
	return out
}

func (s *SystemDTO) GetMetrics() (api.Metrics, error) {
	return nil, api.ErrAlreadyExists
	// return metricsEntitiesToApi(s.service.FetchMetrics()), nil
}

func (s *SystemDTO) LogMetrics(metrics api.Metrics) error {
	return s.service.LogMetrics(metricsApiToEntities(metrics)...)
}

func atoi[Int constraints.Integer](s string) Int {
	i, _ := strconv.Atoi(s)
	return Int(i)
}

func parseBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
