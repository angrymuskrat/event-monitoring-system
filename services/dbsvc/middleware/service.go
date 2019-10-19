package middleware

import (
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbservice"
)

// Middleware describes a service middleware.
type Middleware func(service dbservice.Service) dbservice.Service
