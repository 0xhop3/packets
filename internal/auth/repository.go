package auth

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CREATE_USER = `INSERT INTO users (username, password_hash, public_key, encrypted_privkey, privkey_salt)
	VALUES ($1, $2, $3, $4, $5)`

	GET_BY_USERNAME = `SELECT id, username, password_hash, public_key, encryptedPrivKey, privKey_salt FROM users WHERE username = $1`

	GET_BY_ID = `SELECT id, username, password_hash, public_key, encryptedPricKey, privKey_Salt FROM users WHERE id = $1`

	GET_PUBLIC_KEY = `SELECT public_key FROM users WHERE id = $1`
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx, CREATE_USER, u.Username, u.PasswordHash, u.PublicKey, u.EncryptedPrivKey, u.PrivKeySalt)
	return err
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx, GET_BY_USERNAME, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PublicKey, &u.EncryptedPrivKey, &u.PrivKeySalt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx, GET_BY_ID, id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PublicKey, &u.EncryptedPrivKey, &u.PrivKeySalt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) GetPublicKey(ctx context.Context, userID string) ([]byte, error) {
	var pk []byte
	err := r.db.QueryRow(ctx, GET_PUBLIC_KEY, userID).Scan(&pk)
	return pk, err
}
