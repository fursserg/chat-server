package repository

import (
	"context"
)

type ChatRepository interface {
	Create(ctx context.Context, userIds []int64, name string) (int64, error)
	Delete(ctx context.Context, id int64) error
}
