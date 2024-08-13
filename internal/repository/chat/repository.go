package chat

import (
	"context"
	"encoding/json"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/fursserg/chat-server/internal/client/db"
	"github.com/fursserg/chat-server/internal/repository"
)

const (
	tableName       = "chats"
	idColumn        = "id"
	titleColumn     = "title"
	userIDsColumn   = "user_ids"
	statusColumn    = "status"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
	activeStatus    = 1
	deletedStatus   = 10
)

type repo struct {
	db db.Client
}

// NewRepository Создает репозиторий
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

// Create Создает новый чат
func (r *repo) Create(ctx context.Context, userIds []int64, name string) (int64, error) {
	userIDs, err := json.Marshal(userIds)
	if err != nil {
		log.Fatalf("wrong user_ids: %+v", userIds)
	}

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn, userIDsColumn, statusColumn).
		Values(name, userIDs, activeStatus).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		log.Fatalf("failed to insert chat: %v", err)
	}

	return chatID, nil
}

// Delete Переводит чат в статус "удаленный"
func (r *repo) Delete(ctx context.Context, id int64) error {
	// Вместо удаления, переводим в специальный статус (храним в БД для истории)
	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(statusColumn, deletedStatus).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{"id": id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return nil
}
