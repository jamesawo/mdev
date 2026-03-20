package services

import "context"

type Status struct {
	Running bool
	Message string
}

type Service interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) Status
}
