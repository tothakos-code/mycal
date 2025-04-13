package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang-postgresql-auth-template/internal/models"
)

type EventRepo struct {
	db *sql.DB
}

func NewEventRepo(db *sql.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (e *EventRepo) CreateEvent(ctx context.Context, event models.Event) error {
	query := `INSERT INTO event (user_id, title, description, location, start, finish, is_public, notify_before)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := e.db.ExecContext(ctx, query,
		event.UserID, event.Title, event.Description,
		event.Location, event.Start, event.Finish,
		event.Public, event.NotifyBefore)

	return err
}

func (e *EventRepo) UpdateEvent(ctx context.Context, event models.Event) error {
	query := `UPDATE event
			  SET title = $1, description = $2, location = $3, start = $4, finish = $5, is_public = $6, notify_before = $7
			  WHERE id = $8 AND user_id = $9`

	_, err := e.db.ExecContext(ctx, query,
		event.Title, event.Description, event.Location,
		event.Start, event.Finish, event.Public, event.NotifyBefore,
		event.ID, event.UserID)

	return err
}

func (e *EventRepo) DeleteEvent(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM event WHERE id = $1 AND user_id = $2`

	_, err := e.db.ExecContext(ctx, query, eventID, userID)

	return err
}

func (e *EventRepo) GetEventByID(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	var event models.Event
	query := `SELECT * FROM event WHERE id = $1`
	err := e.db.QueryRowContext(ctx, query, eventID.String()).Scan(
		&event.ID, &event.UserID, &event.Title, &event.Description,
		&event.Location, &event.Start, &event.Finish,
		&event.NotifyBefore, &event.Public, &event.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &event, nil
}

func (e *EventRepo) ListPublicEvents(ctx context.Context) ([]models.Event, error) {
	query := `SELECT * FROM event
			  WHERE is_public = TRUE`

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.UserID,
			&event.Title, &event.Description, &event.Location,
			&event.Start, &event.Finish, &event.NotifyBefore,
			&event.Public, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if events == nil {
		events = []models.Event{}
	}
	return events, nil
}

func (e *EventRepo) ListEventsByUserID(ctx context.Context, userID uuid.UUID, is_public bool) ([]models.Event, error) {
	query := `SELECT * FROM event WHERE user_id = $1 AND is_public = $2`

	rows, err := e.db.QueryContext(ctx, query, userID, is_public)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.UserID,
			&event.Title, &event.Description, &event.Location,
			&event.Start, &event.Finish, &event.NotifyBefore,
			&event.Public, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
