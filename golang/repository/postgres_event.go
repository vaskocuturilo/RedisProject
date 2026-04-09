package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang/domain"
)

type PostgresEventRepository struct {
	db *sql.DB
}

func NewPostgresEventRepository(db *sql.DB) *PostgresEventRepository {
	return &PostgresEventRepository{db: db}
}

func (r *PostgresEventRepository) Create(ctx context.Context, event *domain.Event) error {
	query := `
        INSERT INTO events (id, title, description)
        VALUES ($1, $2, $3)
    `
	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.Title, event.Description)
	return err
}

func (r *PostgresEventRepository) Get(ctx context.Context, id string) (*domain.Event, error) {
	query := `SELECT id, title, description 
              FROM events WHERE id = $1`

	event := &domain.Event{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID, &event.Title, &event.Description)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	return event, err
}

func (r *PostgresEventRepository) GetAll(ctx context.Context) ([]*domain.Event, error) {
	events := make([]*domain.Event, 0)

	query := `SELECT * FROM events`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		e := &domain.Event{}
		if err := rows.Scan(&e.ID, &e.Title, &e.Description); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, err
}

func (r *PostgresEventRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `UPDATE events SET title=$2, description=$3 WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, event.ID, event.Title, event.Description)
	return err
}

func (r *PostgresEventRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return domain.ErrNotFound
	}

	return err
}
