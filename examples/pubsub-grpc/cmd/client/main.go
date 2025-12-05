package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	pb "pubsub-grpc/pkg/pb/pubsub"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port string

func init() {
	flag.StringVar(&port, "port", "12345", "The port to connect to")
	flag.Parse()
}

func main() {
	conn, err := grpc.NewClient(fmt.Sprintf(":%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPubsubServiceClient(conn)

	TryConsumeExample(client)
	time.Sleep(1 * time.Second)
	KeepConsumeExample(client, 5)
}

func PrintMessage(message *pb.Message) {
	fmt.Println("ID:", message.Id)
	fmt.Println("Data:", string(message.Data))
	fmt.Println("Topics:", message.Topics)
	fmt.Println("Publisher:", message.Publisher)
	fmt.Println("CreatedAt:", time.Unix(message.CreatedAt, 0).Format("2006-01-02 15:04:05"))
}
