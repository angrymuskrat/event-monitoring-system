package middleware

import "github.com/angrymuskrat/event-monitoring-system/services/dbconnector"

// Middleware describes a service middleware.
type Middleware func(service dbconnector.Service) dbconnector.Service
