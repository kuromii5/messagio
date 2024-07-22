package transport

import (
	"context"

	"github.com/kuromii5/messagio/internal/models"
	messagio "github.com/kuromii5/messagio/pkg/messagio/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type api struct {
	messagio.UnimplementedMessageServiceServer
	messager Messager
}

type Messager interface {
	SendMessage(ctx context.Context, msg string) (models.SendMessageResponse, error)
	GetMessages(ctx context.Context) (models.GetMessagesResponse, error)
	GetStats(ctx context.Context) (models.GetStatsResponse, error)
}

func NewGrpcServer(messager Messager) *grpc.Server {
	api := &api{messager: messager}

	grpc := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	reflection.Register(grpc)
	messagio.RegisterMessageServiceServer(grpc, api)

	return grpc
}

func (a *api) SendMessage(ctx context.Context, req *messagio.SendMessageRequest) (*messagio.SendMessageResponse, error) {
	response, err := a.messager.SendMessage(ctx, req.Message)
	if err != nil {
		return nil, err
	}

	return &messagio.SendMessageResponse{
		StatusCode: response.StatusCode,
		Message:    response.Message,
	}, nil
}

func (a *api) GetMessages(ctx context.Context, req *messagio.GetMessagesRequest) (*messagio.GetMessagesResponse, error) {
	response, err := a.messager.GetMessages(ctx)
	if err != nil {
		return nil, err
	}

	return &messagio.GetMessagesResponse{Messages: response.Messages}, nil
}

func (a *api) GetStats(ctx context.Context, req *messagio.GetStatsRequest) (*messagio.GetStatsResponse, error) {
	response, err := a.messager.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &messagio.GetStatsResponse{ProcessedCount: response.ProcessedCount}, nil
}
