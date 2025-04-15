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

func (e *EventRepo) ListPublicEvents(ctx context.Context) ([]models.EventWithUser, error) {
	query := `SELECT e.id, e.title, e.description,
		e.location, e.start, e.finish,
		e.notify_before, e.is_public, e.created_at,
		u.id, u.email, u.username,
		u.firstname, u.surname, u.created_at
		FROM event e
		JOIN "users" u ON u.id = e.user_id
		WHERE is_public = TRUE`

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.EventWithUser
	for rows.Next() {
		var event models.EventWithUser
		if err := rows.Scan(&event.ID,
			&event.Title, &event.Description, &event.Location,
			&event.Start, &event.Finish, &event.NotifyBefore,
			&event.Public, &event.CreatedAt,
			&event.User.ID, &event.User.Email, &event.User.Username,
			&event.User.FirstName, &event.User.SurName, &event.User.CreatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if events == nil {
		events = []models.EventWithUser{}
	}
	return events, nil
}

func (e *EventRepo) ListEventsByUserID(ctx context.Context, userID uuid.UUID, is_public bool) ([]models.EventWithUser, error) {
	query := `SELECT e.id, e.title, e.description,
		e.location, e.start, e.finish,
		e.notify_before, e.is_public, e.created_at,
		u.id, u.email, u.username,
		u.firstname, u.surname, u.created_at
		FROM event e
		JOIN "users" u ON u.id = e.user_id
		WHERE user_id = $1 AND is_public = $2`

	rows, err := e.db.QueryContext(ctx, query, userID, is_public)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.EventWithUser
	for rows.Next() {
		var event models.EventWithUser
		if err := rows.Scan(&event.ID,
			&event.Title, &event.Description, &event.Location,
			&event.Start, &event.Finish, &event.NotifyBefore,
			&event.Public, &event.CreatedAt,
			&event.User.ID, &event.User.Email, &event.User.Username,
			&event.User.FirstName, &event.User.SurName, &event.User.CreatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
