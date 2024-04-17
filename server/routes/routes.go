package routes

import "ypeskov/go_hillel_9/internal/logger"

type Routes struct {
	Log *logger.Logger
}

func New(log *logger.Logger) *Routes {
	return &Routes{
		Log: log,
	}
}
