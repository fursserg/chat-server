package chat

import (
	"context"
	"log"

	"github.com/fursserg/chat-server/internal/client/db"
	"github.com/fursserg/chat-server/internal/repository"
	"github.com/fursserg/chat-server/internal/service"
)

type serv struct {
	repo      repository.ChatRepository
	txManager db.TxManager
}

// NewService Создает объект service
func NewService(repo repository.ChatRepository, txManager db.TxManager) service.ChatService {
	return &serv{repo: repo, txManager: txManager}
}

// Create Создает новый чат
func (s *serv) Create(ctx context.Context, userIds []int64, name string) (int64, error) {
	id, err := s.repo.Create(ctx, userIds, name)

	if err != nil {
		log.Fatalf("failed to insert chat: %v", err)
	}

	return id, nil
}

// Delete Удаляет чат
func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)

	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	return nil
}
