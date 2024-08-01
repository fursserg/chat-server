package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/fursserg/chat-server/db"
	chat "github.com/fursserg/chat-server/pkg/chat_v1"
)

const (
	grpcPort = 50052
)

type server struct {
	chat.UnimplementedChatV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chat.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Create Создает новый чат
func (s *server) Create(ctx context.Context, req *chat.CreateRequest) (*chat.CreateResponse, error) {
	userIDs, err := json.Marshal(req.GetUserIds())
	if err != nil {
		log.Fatalf("wrong user_ids: %+v", req.GetUserIds())
	}

	builderInsert := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "user_ids", "status").
		Values(req.GetName(), userIDs, db.ActiveStatus).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var chatID int64
	err = db.GetConnect().QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Fatalf("failed to insert chat: %v", err)
	}

	return &chat.CreateResponse{
		Id: chatID,
	}, nil
}

// Delete Переводит чат в статус "удаленный"
func (s *server) Delete(ctx context.Context, req *chat.DeleteRequest) (*emptypb.Empty, error) {
	// Вместо удаления, переводим в специальный статус (храним в БД для истории)
	builderUpdate := sq.Update("chats").
		PlaceholderFormat(sq.Dollar).
		Set("status", db.DeletedStatus).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := db.GetConnect().Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update chat: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return new(emptypb.Empty), nil
}

// SendMessage Отправляет сообщение в чат
func (s *server) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Send message data: %+v", req)

	return new(emptypb.Empty), nil
}
