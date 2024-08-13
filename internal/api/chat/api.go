package chat

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/fursserg/chat-server/internal/service"
	desc "github.com/fursserg/chat-server/pkg/chat_v1"
)

type ChatApi struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewApi Создает новую структуру ChatApi
func NewApi(chatService service.ChatService) *ChatApi {
	return &ChatApi{
		chatService: chatService,
	}
}

// Create Создание чата
func (c *ChatApi) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := c.chatService.Create(ctx, req.GetUserIds(), req.GetName())
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

// Delete Удаление чата
func (c *ChatApi) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := c.chatService.Delete(ctx, req.GetId())
	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	return new(emptypb.Empty), nil
}
