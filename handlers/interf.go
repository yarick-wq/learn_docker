package handlers

import "context"

type EmployeeStorage interface {
	Create(ctx context.Context, emp Employee) (Employee, error)
	Collect(ctx context.Context) ([]Employee, error)
	Delete(ctx context.Context, id int) error
}
type Employee struct {
	ID       int
	FullName string
	Position string
}
