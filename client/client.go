package client 

import (
	"log"

	pb "github.com/SendHive/worker-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client pb.TaskServiceClient

func InitClient() pb.TaskServiceClient {
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	Client = pb.NewTaskServiceClient(conn)
	log.Println("gRPC client initialized successfully")
	return Client
}
