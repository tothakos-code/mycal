package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-postgresql-auth-template/internal/models"
  "github.com/lib/pq"
)

type EventRepo struct {
	db *sql.DB
}

func NewEventRepo(db *sql.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (e *EventRepo) CreateEvent(ctx context.Context, event models.Event) error {
	query := `INSERT INTO event (id, calendar_id, description, location, start_date, end_date, start_time, end_time, notify_before)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := e.db.ExecContext(ctx, query,
		event.ID, event.CalendarID, event.Description, event.Location, event.StartDate, event.EndDate,
		event.StartTime, event.EndTime, event.NotifyBefore)
	return err
}

func (e *EventRepo) GetEventByID(ctx context.Context, eventID string) (*models.Event, error) {
	var event models.Event
	query := `SELECT * FROM event WHERE id = $1`
	err := e.db.QueryRowContext(ctx, query, eventID).Scan(
		&event.ID, &event.CalendarID, &event.Description, &event.Location,
		&event.StartDate, &event.EndDate, &event.StartTime, &event.EndTime,
		&event.NotifyBefore, &event.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &event, nil
}

func (e *EventRepo) ListPublicEvents(ctx context.Context) ([]models.Event, error) {
	query := `SELECT e.* FROM event e
			  JOIN calendar c ON e.calendar_id = c.id
			  WHERE c.is_public = TRUE`

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.CalendarID, &event.Description,
			&event.Location, &event.StartDate, &event.EndDate, &event.StartTime,
			&event.EndTime, &event.NotifyBefore, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (e *EventRepo) ListEventsByCalendars(ctx context.Context, calendarIDs []string) ([]models.Event, error) {
	query := `SELECT * FROM event WHERE calendar_id = ANY($1)`

	rows, err := e.db.QueryContext(ctx, query, pq.Array(calendarIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.CalendarID, &event.Description,
			&event.Location, &event.StartDate, &event.EndDate, &event.StartTime,
			&event.EndTime, &event.NotifyBefore, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
