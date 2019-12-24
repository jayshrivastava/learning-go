package main

import (
	"context"
	// "encoding/json"
	"flag"
	"fmt"
	// "io"
	// "io/ioutil"
	"log"
	// "math"
	"net"
	// "sync"
	// "time"

	"google.golang.org/grpc"

	// "google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/testdata"

	// "github.com/golang/protobuf/proto"

	pb "../chat"
)

var (
	port       = flag.Int("port", 10000, "The server port")
)

type chatServer struct {
    pb.UnimplementedChatServer
    items map[int32]string
}

func (s *chatServer) GetItem(ctx context.Context, key *pb.ItemKey) (*pb.ItemValue, error) {
    value, found := s.items[key.Index]
    if found {
        return &pb.ItemValue{Value: value}, nil
    }
	return &pb.ItemValue{Value: ""}, nil
}

// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *chatServer) ListItems(rng *pb.Range, stream pb.Chat_ListItemsServer) error {
	for key := rng.StartIndex; key < rng.EndIndex; key++ {
        value, found := s.items[key]
        if found {
            if err := stream.Send(&pb.ItemValue{Value: value}); err != nil {
                return err
            }
        }
	}
	return nil
}

func makeNewServer() *chatServer {
    server := chatServer{items: make(map[int32]string)}
    server.seed()
    return &server
}

func (server *chatServer) seed() {
    server.items[1] = "random"
    server.items[2] = "hello"
    server.items[3] = "google"
    server.items[4] = "golang"
    server.items[5] = "computer"
    server.items[6] = "semaphore"
    server.items[7] = "pi"
    server.items[8] = "example"
    server.items[9] = "integral"
    server.items[10] = "gRPC"
}

func main() {
    flag.Parse() // parse flag variables 

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()

	pb.RegisterChatServer(grpcServer, makeNewServer())
	grpcServer.Serve(lis)
}