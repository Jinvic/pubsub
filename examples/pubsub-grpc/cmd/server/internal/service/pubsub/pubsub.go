package pubsub

import (
	pb "pubsub-grpc/pkg/pb/pubsub"
)

type PubsubService struct {
	pb.UnimplementedPubsubServiceServer
}
