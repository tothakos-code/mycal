package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-postgresql-auth-template/internal/models"
  "github.com/google/uuid"
)

type CalendarRepo struct {
	db *sql.DB
}

func NewCalendarRepo(db *sql.DB) *CalendarRepo {
	return &CalendarRepo{db: db}
}

func (c *CalendarRepo) CreateCalendar(ctx context.Context, userID uuid.UUID, isPublic bool) error {
	query := `INSERT INTO calendar (user_id, is_public) VALUES ($1, $2)`
	_, err := c.db.ExecContext(ctx, query, userID, isPublic)
	return err
}

func (c *CalendarRepo) GetCalendarsByUserID(ctx context.Context, userID string) ([]models.Calendar, error) {
	query := `SELECT * FROM calendar WHERE user_id = $1`
	rows, err := c.db.QueryContext(ctx, query, userID)
  if err != nil {
    return nil, fmt.Errorf("database error: %w", err)
  }
	defer rows.Close()

	var calendars []models.Calendar
	for rows.Next() {
		var cal models.Calendar
		if err := rows.Scan(&cal.ID, &cal.UserID, &cal.IsPublic, &cal.CreatedAt); err != nil {
			return nil, err
		}
		calendars = append(calendars, cal)
	}
	return calendars, nil
}

func (c *CalendarRepo) GetPublicCalendar(ctx context.Context) (models.Calendar, error) {
	query := `SELECT * FROM calendar WHERE is_public = TRUE`
	rows, err := c.db.QueryContext(ctx, query)
  if err != nil {
    return models.Calendar{}, fmt.Errorf("database error: %w", err)
  }
	defer rows.Close()

	var calendar models.Calendar
  if err := rows.Scan(&calendar.ID, &calendar.UserID, &calendar.IsPublic, &calendar.CreatedAt); err != nil {
    return models.Calendar{}, err
  }
	return calendar, nil
}
