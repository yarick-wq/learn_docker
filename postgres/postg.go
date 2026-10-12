package postgres

import (
	"context"
	"docker/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) Create(ctx context.Context, emp handlers.Employee) (handlers.Employee, error) {
	sqlquery := `
	INSERT INTO employees(full_name,position)
	VALUES($1,$2)
	RETURNING id,full_name,position;
	`
	row := s.pool.QueryRow(ctx, sqlquery, emp.FullName, emp.Position)
	if err := row.Scan(&emp.ID, &emp.FullName, &emp.Position); err != nil {
		return handlers.Employee{}, err
	}
	return emp, nil
}

func (s *Storage) Collect(ctx context.Context) ([]handlers.Employee, error) {
	sqlquery := `
	SELECT id,full_name,position FROM employees
	ORDER BY id ASC
	`

	rows, err := s.pool.Query(ctx, sqlquery)
	if err != nil {
		return nil, err
	}
	var Employee handlers.Employee
	defer rows.Close()
	emp := make([]handlers.Employee, 0)
	for rows.Next() {
		if err := rows.Scan(
			&Employee.ID,
			&Employee.FullName,
			&Employee.Position,
		); err != nil {
			return emp, err
		}
		emp = append(emp, Employee)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	sqlquery := `
	DELETE FROM employees
	WHERE id=$1
	`
	tag, err := s.pool.Exec(ctx, sqlquery, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return handlers.ErrIdNotFound
	}
	return nil
}
