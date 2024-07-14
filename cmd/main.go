package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	chat "github.com/fursserg/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

type server struct {
	chat.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *chat.CreateRequest) (*chat.CreateResponse, error) {
	log.Printf("Chat create data: %+v", req.GetUsernames())

	return &chat.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Delete(ctx context.Context, req *chat.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Chat delete id: %d", req.GetId())

	res := new(emptypb.Empty)
	return res, nil
}

func (s *server) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Send message data: %+v", req)

	res := new(emptypb.Empty)
	return res, nil
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
