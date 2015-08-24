package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"google.golang.org/grpc/grpclog"

	pb "github.com/gotalk2/proto"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type server struct {
}

func (*server) GetHelloWorld(context context.Context, clientMsg *pb.Msg) (*pb.Msg, error) {
	content := "hello " + clientMsg.Content
	msg := &pb.Msg{Content: content}
	return msg, nil
}

// 	LetsStreaming(TalkMessage_LetsStreamingServer) error
func (*server) LetsStreaming(stream pb.TalkMessage_LetsStreamingServer) error {
	for {
		in, err := stream.Recv()
		// end of the streaming
		if err == io.EOF {
			grpclog.Println("finished stream")
			return nil
		}
		if err != nil {
			grpclog.Printf("returned with error %v", err)
			return err
		}
		content := in.Content
		revMsg := "received message from streaming: " + content + "->" + strings.ToUpper(content)
		grpclog.Println(revMsg)
		sleep := 5
		grpclog.Printf("wait for %v seconds for response", sleep)
		time.Sleep(time.Duration(sleep) * time.Second) //"sleep for 5 seconds"

		stream.Send(&pb.Msg{Content: revMsg})
	}

}

func newServer() *server {
	s := new(server)
	return s
}

func main() {
	grpclog.Println("start server...")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatal("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTalkMessageServer(grpcServer, newServer())
	grpcServer.Serve(lis)
	grpclog.Println("server shutdown...")

}
