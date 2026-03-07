package contact

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddContact(ctx context.Context, userID, contactID string) error {
	if userID == contactID {
		return fmt.Errorf("cannot add yourself as a contact")
	}

	return s.repo.Add(ctx, userID, contactID)
}

func (s *Service) RemoveContact(ctx context.Context, userID, contactID string) error {
	return s.repo.Remove(ctx, userID, contactID)
}

func (s *Service) BlockContact(ctx context.Context, userID, contactID string) error {
	return s.repo.Block(ctx, userID, contactID)
}

func (s *Service) ListContacts(ctx context.Context, userID string) ([]Contact, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) IsBlocked(ctx context.Context, senderID, recipientID string) (bool, error) {
	return s.repo.IsBlocked(ctx, senderID, recipientID)
}
