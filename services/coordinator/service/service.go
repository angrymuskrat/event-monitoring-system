package service

import "github.com/angrymuskrat/event-monitoring-system/services/coordinator/service/status"

type CoordinatorService interface {
	NewSession(req Session) (string, error)
	Status(id string) (status.Status, error)
}
