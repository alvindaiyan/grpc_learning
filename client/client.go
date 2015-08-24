package main

import (
	"flag"
	"io"
	"time"

	pb "github.com/gotalk2/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"

	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
)

func getHelloWorld(client pb.TalkMessageClient, content string) {
	msg, err := client.GetHelloWorld(context.Background(), &pb.Msg{Content: content})
	if err != nil {
		grpclog.Fatalf("fail to call getHelloWorld: %v", err)
	}

	grpclog.Println("get Message from Server: ", msg.Content)
}

func letsStreaming(client pb.TalkMessageClient) {
	// messages
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h"}

	stream, err := client.LetsStreaming(context.Background())
	if err != nil {
		grpclog.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}

	// receive msg
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done
				close(waitc)
				return
			}
			if err != nil {
				grpclog.Fatalf("Failed to receive a note : %v", err)
			}
			grpclog.Printf("Got message %s", in.Content)
		}
	}()

	// send msg
	for _, str := range letters {
		grpclog.Printf("send msg: %v", str)
		if err := stream.Send(&pb.Msg{Content: str}); err != nil {
			grpclog.Fatalf("%v.Send(%v) = %v", stream, str, err)
		}
		sleep := 5
		grpclog.Printf("sleep for %v seconds", sleep)
		time.Sleep(time.Duration(sleep) * time.Second) //"sleep for 5 seconds"
	}

	// close send
	stream.CloseSend()
	<-waitc
}

func main() {
	conn, err := grpc.Dial(*serverAddr)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	// always turn off the conn when exit
	defer conn.Close()

	client := pb.NewTalkMessageClient(conn)

	// getHelloWorld(client, "from client")
	letsStreaming(client)

}
