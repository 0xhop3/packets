package presence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb *redis.Client
}

func presenceKey(userID string) string { return fmt.Sprintf("presence:%s", userID) }
func channelKey(userID string) string  { return fmt.Sprintf("notify:%s", userID) }

func (s *Service) SetOnline(ctx context.Context, userID string) error {
	return s.rdb.Set(ctx, presenceKey(userID), "1", 60*time.Second).Err()
}

func (s *Service) SetOffline(ctx context.Context, userID string) error {
	return s.rdb.Del(ctx, presenceKey(userID)).Err()
}

func (s *Service) IsOnline(ctx context.Context, userID string) (bool, error) {
	res, err := s.rdb.Exists(ctx, presenceKey(userID)).Result()
	return res > 0, err
}

func (s *Service) HeartBeat(ctx context.Context, userID string) error {
	return s.rdb.Expire(ctx, presenceKey(userID), 60*time.Second).Err()
}

type Notification struct {
	Type     string `json:"type"`
	FromID   string `json:"from_id"`
	FromUser string `json:"from_user"`
	Preview  string `json:"preview"`
}

func (s *Service) Publish(ctx context.Context, recipientID string, notify Notification) error {
	data, _ := json.Marshal(notify)
	return s.rdb.Publish(ctx, channelKey(recipientID), data).Err()
}

func (s *Service) Subscribe(ctx context.Context, userID string) (<-chan Notification, func()) {
	sub := s.rdb.Subscribe(ctx, channelKey(userID))
	ch := make(chan Notification, 32)

	go func() {
		defer close(ch)
		for msg := range sub.Channel() {
			var n Notification
			if json.Unmarshal([]byte(msg.Payload), &n) == nil {
				ch <- n
			}
		}
	}()

	cancel := func() { sub.Close() }
	return ch, cancel
}
