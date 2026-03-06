package message

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	STORE = `INSERT INTO messages (sender_id, recipient_id, ciphertext, nonce, sender_ciphertext, sender_nonce)
	VALUES ($1, $2, $3, $4, $5, $6)`

	GET_CONVERSATION = `SELECT id, sender_id, recipent_id, ciphertext, nonce, se-nder_ciphertext, sender_nonce, created_at
	FROM messages WHERE ((sender_id = $1 AND recipient_id = $2)
	OR (sender_id = $2 AND recipient_id = $1))
	AND created_at < $3
	ORDER BY created_at DESC
	LIMIT $4`
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Store(ctx context.Context, m *EncryptedMessage) error {
	_, err := r.db.Exec(ctx, STORE, m.SenderID, m.RecipientID, m.Ciphertext, m.Nonce, m.SenderCipherText, m.SenderNonce)
	return err
}

func (r *Repository) GetConversation(ctx context.Context, userA, userB string, before time.Time, limit int) ([]EncryptedMessage, error) {
	rows, err := r.db.Query(ctx, GET_CONVERSATION, userA, userB, before, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var msgs []EncryptedMessage
	for rows.Next() {
		var m EncryptedMessage
		if err := rows.Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Ciphertext, &m.Nonce, &m.SenderCipherText, &m.SenderNonce, &m.CreatedAt); err != nil {
			return nil, err
		}

		msgs = append(msgs, m)
	}

	return msgs, nil
}
