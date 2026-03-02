package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/0xhop3/packets/internal/crypto"
	"golang.org/x/crypto/argon2"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	return fmt.Sprintf("$argon2id$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func verifyPassword(password, encoded string) bool {
	parts := strings.SplitN(encoded, "$", 4)
	if len(parts) != 4 || parts[1] != "argon2id" {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}

	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}

func (s *Service) Register(ctx context.Context, username, password string) (*SessionIdentity, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash: %w", err)
	}

	kp, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("keygen: %w", err)
	}

	encPriv, nonce, salt, err := crypto.EncryptPrivateKey(kp.PrivateKey, password)
	if err != nil {
		return nil, fmt.Errorf("encrypt privkey: %w", err)
	}

	encWithNonce := append(nonce, encPriv...)

	u := &User{
		Username:         username,
		PasswordHash:     passwordHash,
		PublicKey:        kp.PublicKey[:],
		EncryptedPrivKey: encWithNonce,
		PrivKeySalt:      salt,
	}

	if err := s.repo.CreateUser(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	created, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &SessionIdentity{
		UserID:     created.ID,
		Username:   username,
		PublicKey:  kp.PublicKey,
		PrivateKey: kp.PrivateKey,
	}, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (*SessionIdentity, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !verifyPassword(password, u.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if len(u.EncryptedPrivKey) < 12 {
		return nil, errors.New("corrupt key data")
	}

	nonce := u.EncryptedPrivKey[:12]
	encPriv := u.EncryptedPrivKey[12:]

	privKey, err := crypto.DecryptPrivateKey(encPriv, nonce, u.PrivKeySalt, password)
	if err != nil {
		return nil, fmt.Errorf("decrypt privkey: %w", err)
	}

	var pubKey [32]byte
	copy(pubKey[:], u.PublicKey)

	return &SessionIdentity{
		UserID:     u.ID,
		Username:   u.Username,
		PublicKey:  pubKey,
		PrivateKey: privKey,
	}, nil
}

func (s *Service) GetPublicKey(ctx context.Context, userID string) ([32]byte, error) {
	pk, err := s.repo.GetPublicKey(ctx, userID)
	if err != nil {
		return [32]byte{}, err
	}

	var key [32]byte
	copy(key[:], pk)
	return key, nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.GetByUsername(ctx, username)
}
