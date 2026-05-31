package main

import (
	"context"
	"testing"

	pb "github.com/TerrenceHo/monorepo/proto/api/v1/greeter"
)

func TestSayHello_Default(t *testing.T) {
	s := &server{}
	resp, err := s.SayHello(context.Background(), &pb.HelloRequest{
		Name: "World",
	})
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	want := "Hello, World!"
	if resp.GetMessage() != want {
		t.Errorf("got %q, want %q", resp.GetMessage(), want)
	}
	if resp.GetTimestampUnix() == 0 {
		t.Error("expected non-zero timestamp")
	}
}

func TestSayHello_Formal(t *testing.T) {
	s := &server{}
	resp, err := s.SayHello(context.Background(), &pb.HelloRequest{
		Name:  "Alice",
		Style: pb.GreetingStyle_GREETING_STYLE_FORMAL,
	})
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	want := "Good day, Alice. It is a pleasure to make your acquaintance."
	if resp.GetMessage() != want {
		t.Errorf("got %q, want %q", resp.GetMessage(), want)
	}
}
