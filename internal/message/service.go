package message

import (
	"context"
	"fmt"
	"time"

	"github.com/0xhop3/packets/internal/crypto"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendEncrypted(ctx context.Context, senderID, recipientID string, plaintext []byte, senderPriv, recipientPub, senderPub [32]byte) error {
	recipienShared, err := crypto.ECDH(senderPriv, recipientPub)
	if err != nil {
		return fmt.Errorf("ecdh recipient: %w", err)
	}

	recipientAESKey, err := crypto.DeriveAESKey(recipienShared, []byte("packets-msg-v1"))
	if err != nil {
		return fmt.Errorf("derive key: %w", err)
	}

	ct, nonce, err := crypto.Encrypt(recipientAESKey, plaintext)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	senderShared, err := crypto.ECDH(senderPriv, senderPub)
	if err != nil {
		return fmt.Errorf("ecdh sender: %w", err)
	}

	senderAESKey, err := crypto.DeriveAESKey(senderShared, []byte("packets-msg-v1-self"))
	if err != nil {
		return err
	}

	senderCT, senderNonce, err := crypto.Encrypt(senderAESKey, plaintext)
	if err != nil {
		return err
	}

	return s.repo.Store(ctx, &EncryptedMessage{
		SenderID:         senderID,
		RecipientID:      recipientID,
		Ciphertext:       ct,
		Nonce:            nonce,
		SenderCipherText: senderCT,
		SenderNonce:      senderNonce,
	})
}

func (s *Service) DecryptMessage(msg EncryptedMessage, readerID string, readerPriv, peerPub, readerPub [32]byte) ([]byte, error) {
	if readerID == msg.SenderID {
		shared, err := crypto.ECDH(readerPriv, readerPub)
		if err != nil {
			return nil, err
		}

		key, err := crypto.DeriveAESKey(shared, []byte("packets-msg-v1-self"))
		if err != nil {
			return nil, err
		}

		return crypto.Decrypt(key, msg.SenderCipherText, msg.SenderNonce)
	}

	shared, err := crypto.ECDH(readerPriv, peerPub)
	if err != nil {
		return nil, err
	}

	key, err := crypto.DeriveAESKey(shared, []byte("packets-msg-v1"))
	if err != nil {
		return nil, err
	}

	return crypto.Decrypt(key, msg.Ciphertext, msg.Nonce)
}

func (s *Service) GetConversation(ctx context.Context, userA, userB string, before time.Time, limit int) ([]EncryptedMessage, error) {
	return s.repo.GetConversation(ctx, userA, userB, before, limit)
}
