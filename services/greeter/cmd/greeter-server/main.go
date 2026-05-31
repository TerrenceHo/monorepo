package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/TerrenceHo/monorepo/proto/api/v1/greeter"
)

var port = flag.Int("port", 50051, "The server port")

// server implements the GreeterService gRPC interface.
type server struct {
	// Embedding the unimplemented server ensures forward compatibility:
	// if new RPCs are added to the proto, this still compiles.
	pb.UnimplementedGreeterServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	greeting := formatGreeting(req.GetName(), req.GetStyle())
	log.Printf("SayHello: name=%q style=%v", req.GetName(), req.GetStyle())

	return &pb.HelloReply{
		Message:       greeting,
		TimestampUnix: time.Now().Unix(),
	}, nil
}

func (s *server) SayHelloStream(req *pb.HelloRequest, stream pb.GreeterService_SayHelloStreamServer) error {
	name := req.GetName()
	style := req.GetStyle()
	log.Printf("SayHelloStream: name=%q style=%v", name, style)

	// Send 5 greetings, one per second.
	for i := 0; i < 5; i++ {
		greeting := fmt.Sprintf("[%d/5] %s", i+1, formatGreeting(name, style))
		err := stream.Send(&pb.HelloReply{
			Message:       greeting,
			TimestampUnix: time.Now().Unix(),
		})
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

func formatGreeting(name string, style pb.GreetingStyle) string {
	switch style {
	case pb.GreetingStyle_GREETING_STYLE_FORMAL:
		return fmt.Sprintf("Good day, %s. It is a pleasure to make your acquaintance.", name)
	case pb.GreetingStyle_GREETING_STYLE_CASUAL:
		return fmt.Sprintf("Hey %s, what's up!", name)
	case pb.GreetingStyle_GREETING_STYLE_ENTHUSIASTIC:
		return fmt.Sprintf("OH WOW, %s!!! SO GREAT TO SEE YOU!!!", name)
	default:
		return fmt.Sprintf("Hello, %s!", name)
	}
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServiceServer(s, &server{})

	// Enable server reflection so clients like grpcurl can introspect.
	reflection.Register(s)

	log.Printf("greeter-server listening on :%d", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
