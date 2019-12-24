package main

import (
	"context"
	"flag"
	"log"
	"time"
	"io"

	"google.golang.org/grpc"
	pb "../chat"
)

var (
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
)

func printItem(client pb.ChatClient, key *pb.ItemKey) {
	log.Printf("Getting message at index %d", key.Index)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	message, err := client.GetItem(ctx, key)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	log.Println(message)
}

// printFeatures lists all the features within the given bounding Rectangle.
func printItems(client pb.ChatClient, rng *pb.Range) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListItems(ctx, rng)
	if err != nil {
		log.Fatalf("%v.ListItems(_) = _, %v", client, err)
	}
	for {
		itemValue, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListItems(_) = _, %v", client, err)
		}
		log.Println(itemValue)
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewChatClient(conn)

	// Looking for a valid feature
	printItem(client, &pb.ItemKey{Index: 1})

	//get all items 1 by 1
	printItems(client, &pb.Range{StartIndex: 2, EndIndex: 9})
}