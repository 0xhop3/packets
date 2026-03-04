package contact

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

const (
	ADD = `INSERT INTO contact (user_id, contact_id) VALUES ($1, $2) ON CONFLICT (user_id, contact_id) DO UPDATE
	SET BLOCKED = false`

	REMOVE = `DELETE FROM contacts WHERE user_id = $1 AND contact_id = $2`

	BLOCK = `INSERT INTO contacts (user_id, contact_id, blocked) VALUES ($1, $2, true) ON CONFLICT
	(user_id, contact_id) DO UPDATE SET UPDATE SET blocked = true`

	LIST = `SELECT c.id, c.user_id, c.contact_id, u.username, c.blocked FROM contacts c JOIN users ON u.id = c.contact_id
	WHERE c.user_id = $1 OREDER BY u.username`

	IS_BLOCKED = `SELECT blocked FROM contacts WHERE user_id = $1 AND contact_id = $2`
)

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, userID, contactID string) error {
	_, err := r.db.Exec(ctx, ADD, userID, contactID)
	return err
}

func (r *Repository) Remove(ctx context.Context, userID, contactID string) error {
	_, err := r.db.Exec(ctx, REMOVE, userID, contactID)
	return err
}

func (r *Repository) Block(ctx context.Context, userID, contactID string) error {
	_, err := r.db.Exec(ctx, BLOCK, userID, contactID)
	return err
}

func (r *Repository) List(ctx context.Context, userID string) ([]Contact, error) {
	rows, err := r.db.Query(ctx, LIST, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var c Contact
		if err := rows.Scan(&c.ID, &c.UserID, &c.ContactID, &c.Username, &c.Blocked); err != nil {
			return nil, err
		}

		contacts = append(contacts, c)
	}

	return contacts, nil
}

func (r *Repository) IsBlocked(ctx context.Context, userID, byUserID string) (bool, error) {
	var blocked bool
	err := r.db.QueryRow(ctx, IS_BLOCKED, byUserID, userID).Scan(&blocked)
	if err != nil {
		return false, nil
	}

	return blocked, nil
}
