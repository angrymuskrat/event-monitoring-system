package service

type CoordinatorService interface {
	NewSession(req SessionParameters) (string, error)
	Status(id string) (string, error)
}
