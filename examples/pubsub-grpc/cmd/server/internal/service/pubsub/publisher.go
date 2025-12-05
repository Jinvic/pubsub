package pubsub

import (
	"context"
	pb "pubsub-grpc/pkg/pb/pubsub"

	"github.com/Jinvic/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetPublisher(pubID string) (*pubsub.Publisher, error) {
	if pubID == "" {
		return nil, status.Error(codes.InvalidArgument, "pub_id is required")
	}
	pub := pubsub.GetPublisher(pubsub.PubSubID(pubID))
	if pub == nil {
		return nil, status.Errorf(codes.NotFound, "publisher not found: %s", pubID)
	}
	return pub, nil
}

func (s *PubsubService) NewPublisher(ctx context.Context, empty *emptypb.Empty) (*pb.NewPublisherResponse, error) {
	pub := pubsub.NewPublisher(true)
	resp := &pb.NewPublisherResponse{
		PubId: string(pub.ID),
	}
	return resp, nil
}

func (s *PubsubService) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	pub, err := GetPublisher(req.PubId)
	if err != nil {
		return nil, err
	}
	pub.Publish(req.Topics, req.Data)
	return &emptypb.Empty{}, nil
}

func (s *PubsubService) UnregisterPublisher(ctx context.Context, req *pb.UnregisterPublisherRequest) (*emptypb.Empty, error) {
	pub, err := GetPublisher(req.PubId)
	if err != nil {
		return nil, err
	}
	pub.Unregister()
	return &emptypb.Empty{}, nil
}
