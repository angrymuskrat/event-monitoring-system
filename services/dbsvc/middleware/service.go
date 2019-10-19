package middleware

import "github.com/angrymuskrat/event-monitoring-system/services/dbsvc"

// Middleware describes a service middleware.
type Middleware func(service dbsvc.Service) dbsvc.Service
