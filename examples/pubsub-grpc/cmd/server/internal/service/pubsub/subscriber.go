package pubsub

import (
	"context"
	"errors"

	pb "pubsub-grpc/pkg/pb/pubsub"

	"github.com/Jinvic/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetSubscriber(subID string) (*pubsub.Subscriber, error) {
	if subID == "" {
		return nil, status.Error(codes.InvalidArgument, "sub_id is required")
	}
	sub := pubsub.GetSubscriber(pubsub.PubSubID(subID))
	if sub == nil {
		return nil, status.Errorf(codes.NotFound, "subscriber not found: %s", subID)
	}
	return sub, nil
}

func (s *PubsubService) NewSubscriber(ctx context.Context, req *pb.NewSubscriberRequest) (*pb.NewSubscriberResponse, error) {
	sub := pubsub.NewSubscriber(int(req.BufferSize), true)
	resp := &pb.NewSubscriberResponse{
		SubId: string(sub.ID),
	}
	return resp, nil
}

func (s *PubsubService) UnregisterSubscriber(ctx context.Context, req *pb.UnregisterSubscriberRequest) (*emptypb.Empty, error) {
	sub, err := GetSubscriber(req.SubId)
	if err != nil {
		return nil, err
	}
	sub.Unregister()
	return &emptypb.Empty{}, nil
}

func (s *PubsubService) Subscribe(ctx context.Context, req *pb.SubscribeRequest) (*emptypb.Empty, error) {
	sub, err := GetSubscriber(req.SubId)
	if err != nil {
		return nil, err
	}
	sub.Subscribe(pubsub.PubSubID(req.PubId), req.Topic)
	return &emptypb.Empty{}, nil
}

func (s *PubsubService) UnsubscribePublisher(ctx context.Context, req *pb.UnsubscribePublisherRequest) (*emptypb.Empty, error) {
	sub, err := GetSubscriber(req.SubId)
	if err != nil {
		return nil, err
	}
	sub.UnsubscribePublisher(pubsub.PubSubID(req.PubId))
	return &emptypb.Empty{}, nil
}

func (s *PubsubService) UnsubscribeTopic(ctx context.Context, req *pb.UnsubscribeTopicRequest) (*emptypb.Empty, error) {
	sub, err := GetSubscriber(req.SubId)
	if err != nil {
		return nil, err
	}
	sub.UnsubscribeTopic(pubsub.PubSubID(req.PubId), req.Topics)
	return &emptypb.Empty{}, nil
}

func (s *PubsubService) TryConsume(ctx context.Context, req *pb.TryConsumeRequest) (*pb.TryConsumeResponse, error) {
	sub, err := GetSubscriber(req.SubId)
	if err != nil {
		return nil, err
	}

	// 非阻塞获取消息
	msg, ok, err := sub.TryConsume()
	if err != nil {
		// 比如 subscriber 已关闭
		return nil, status.Errorf(codes.FailedPrecondition, "consume failed: %v", err)
	}

	var message *pb.Message
	if ok {
		message = &pb.Message{
			Id:        string(msg.ID),
			Data:      msg.Data,
			Topics:    msg.Topics,
			Publisher: string(msg.Publisher),
			CreatedAt: msg.CreatedAt.Unix(),
		}
	}
	return &pb.TryConsumeResponse{
		Message: message,
		Ok:      ok,
	}, nil
}

func (s *PubsubService) KeepConsume(req *pb.KeepConsumeRequest, stream pb.PubsubService_KeepConsumeServer) error {
	sub, err := GetSubscriber(req.SubId)
	if err != nil {
		return err
	}

	ctx := stream.Context()
	for {
		// 阻塞获取消息
		msg, err := sub.ConsumeContext(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				// 客户端主动取消或超时，正常退出
				return nil
			}
			if errors.Is(err, pubsub.ErrSubscriberClosed) {
				return status.Error(codes.Canceled, "subscriber closed")
			}
			// 其他未知错误
			return status.Errorf(codes.Internal, "consume error: %v", err)
		}

		message := &pb.Message{
			Id:        string(msg.ID),
			Data:      msg.Data,
			Topics:    msg.Topics,
			Publisher: string(msg.Publisher),
			CreatedAt: msg.CreatedAt.Unix(),
		}
		stream.Send(&pb.KeepConsumeResponse{Message: message})
	}
}
