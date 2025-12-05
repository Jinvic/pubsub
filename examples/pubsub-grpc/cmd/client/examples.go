package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	pb "pubsub-grpc/pkg/pb/pubsub"
)

func TryConsumeExample(client pb.PubsubServiceClient) {
	// 创建 publisher 和 subscriber
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pubID, err := client.NewPublisher(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}
	subID, err := client.NewSubscriber(ctx, &pb.NewSubscriberRequest{
		BufferSize: 10,
	})
	if err != nil {
		log.Fatalf("failed to create subscriber: %v", err)
	}
	_, err = client.Subscribe(ctx, &pb.SubscribeRequest{
		SubId: subID.SubId,
		PubId: pubID.PubId,
		Topic: "try_consume_example",
	})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
		return
	}
	defer func() {
		client.UnregisterPublisher(ctx, &pb.UnregisterPublisherRequest{
			PubId: pubID.PubId,
		})
		client.UnregisterSubscriber(ctx, &pb.UnregisterSubscriberRequest{
			SubId: subID.SubId,
		})
	}()
	fmt.Println("[system]非阻塞消费示例：")
	time.Sleep(1 * time.Second)

	// 第一次消费，没有消息
	resp, err := client.TryConsume(ctx, &pb.TryConsumeRequest{
		SubId: subID.SubId,
	})
	if err != nil {
		log.Fatalf("failed to consume message: %v", err)
	}
	// fmt.Println("第一次消费，没有消息：", message.Message.Data)
	fmt.Println("[subscriber]第一次消费，没有消息：")
	fmt.Println("resp.Message:", resp.Message, "resp.Ok:", resp.Ok)
	time.Sleep(1 * time.Second)

	// 发布消息
	client.Publish(ctx, &pb.PublishRequest{
		PubId:  pubID.PubId,
		Topics: []string{"try_consume_example"},
		Data:   []byte("Hello, world!"),
	})
	fmt.Println("[publisher]发布消息：Hello, world!")
	time.Sleep(1 * time.Second)

	// 第二次消费，有消息
	resp, err = client.TryConsume(ctx, &pb.TryConsumeRequest{
		SubId: subID.SubId,
	})
	if err != nil {
		log.Fatalf("failed to consume message: %v", err)
	}
	fmt.Println("[subscriber]第二次消费，有消息：")
	PrintMessage(resp.Message)
	time.Sleep(1 * time.Second)

}

func KeepConsumeExample(client pb.PubsubServiceClient, count int) {
	// 创建 publisher 和 subscriber
	setupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pubID, err := client.NewPublisher(setupCtx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}
	subID, err := client.NewSubscriber(setupCtx, &pb.NewSubscriberRequest{
		BufferSize: 10,
	})
	if err != nil {
		log.Fatalf("failed to create subscriber: %v", err)
	}
	_, err = client.Subscribe(setupCtx, &pb.SubscribeRequest{
		SubId: subID.SubId,
		PubId: pubID.PubId,
		Topic: "keep_consume_example",
	})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
		return
	}
	defer func() {
		client.UnregisterPublisher(setupCtx, &pb.UnregisterPublisherRequest{
			PubId: pubID.PubId,
		})
		client.UnregisterSubscriber(setupCtx, &pb.UnregisterSubscriberRequest{
			SubId: subID.SubId,
		})
	}()
	fmt.Println("[system]阻塞消费示例：")

	// 随机时间间隔发送消息
	go func() {
		// 发送 count 次消息
		for i := range count {
			time.Sleep(time.Duration(1000+rand.Intn(3000)) * time.Millisecond)
			data := fmt.Sprintf("message %d", i+1)
			client.Publish(context.Background(), &pb.PublishRequest{
				PubId:  pubID.PubId,
				Topics: []string{"keep_consume_example"},
				Data:   []byte(data),
			})
			fmt.Println("[publisher]发布第", i+1, "条消息")
		}
		fmt.Println("[publisher]发送消息完成")
	}()

	// 阻塞消费消息
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stream, err := client.KeepConsume(ctx, &pb.KeepConsumeRequest{
		SubId: subID.SubId,
	})
	if err != nil {
		log.Fatalf("failed to consume message: %v", err)
	}

	for range count {
		resp, err := stream.Recv()
		if err != nil {
			// 检查是否是流正常结束（如服务端关闭）
			if err == io.EOF {
				log.Printf("stream ended early")
				break
			}
			log.Fatalf("failed to consume message: %v", err)
		}
		fmt.Println("[subscriber]消费消息：")
		PrintMessage(resp.Message)
	}
	fmt.Println("[subscriber]消费消息完成")

}
