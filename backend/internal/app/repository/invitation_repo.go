package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-postgresql-auth-template/internal/models"
)

type InvitationRepo struct {
	db *sql.DB
}

func NewInvitationRepo(db *sql.DB) *InvitationRepo {
	return &InvitationRepo{db: db}
}

func (i *InvitationRepo) CreateInvitation(ctx context.Context, userID, eventID, status string) error {
	query := `INSERT INTO invitation (user_id, event_id, status) VALUES ($1, $2, $3)`
	_, err := i.db.ExecContext(ctx, query, userID, eventID, status)
	return err
}

func (i *InvitationRepo) GetInvitationsByUserID(ctx context.Context, userID string) ([]models.Invitation, error) {
	query := `SELECT * FROM invitation WHERE user_id = $1`
	rows, err := i.db.QueryContext(ctx, query, userID)
  if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var invitations []models.Invitation
	for rows.Next() {
		var invite models.Invitation
		if err := rows.Scan(&invite.ID, &invite.UserID, &invite.EventID, &invite.Status, &invite.CreatedAt); err != nil {
			return nil, err
		}
		invitations = append(invitations, invite)
	}
	return invitations, nil
}

func (i *InvitationRepo) ListInvitedUsersByEventID(ctx context.Context, eventID string) ([]models.User, error) {
	query := `SELECT u.id, u.username, u.email, u.first_name, u.surname
			  FROM invitation i
			  JOIN users u ON i.user_id = u.id
			  WHERE i.event_id = $1`

	rows, err := i.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email,
			&user.FirstName, &user.SurName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
