package service

import "github.com/angrymuskrat/event-monitoring-system/services/coordinator/service/status"

type coordinatorService struct {
}

func (s coordinatorService) NewSession(req Session) (string, error) {

}

func (s coordinatorService) Status(id string) (status.Status, error) {

}
