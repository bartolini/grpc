package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/bartolini/grpc/miaow"
	"google.golang.org/grpc"
)

const listenAndServe = ":5000"
const msgCount = 10

type serverImplementation struct {
	msgCount int
	miaow.UnimplementedChatServiceServer
}

func (s *serverImplementation) SayHello(ctx context.Context, msg *miaow.Message) (*miaow.Message, error) {
	s.msgCount++
	log.Printf("Server received [%03d]: %s", s.msgCount, msg.Body)
	reply := new(miaow.Message)
	reply.Body = fmt.Sprintf("=== PONG %s ===", msg.Body)
	return reply, nil
}

func (s *serverImplementation) StreamHello(stream miaow.ChatService_StreamHelloServer) error {
	for {
		if msg, err := stream.Recv(); err == io.EOF {
			log.Print("Goodbye my stream...")
			return nil
		} else if err != nil {
			return err
		} else {
			s.msgCount++
			log.Printf("Stream received [%03d]: %s", s.msgCount, msg.Body)
			reply := new(miaow.Message)
			reply.Body = fmt.Sprintf("=== STREAM-BACK %s ===", msg.Body)
			stream.Send(reply)
		}
	}
}

func main() {

	listener, err := net.Listen("tcp", listenAndServe)
	if err != nil {
		tryClient()
		return
	}

	log.Print("Starting server...")

	grpcServer := grpc.NewServer()

	miaow.RegisterChatServiceServer(grpcServer, new(serverImplementation))

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func tryClient() {
	log.Print("Starting client...")

	conn, err := grpc.Dial(listenAndServe, grpc.WithInsecure())
	if err != nil {
		log.Fatal("unable to connect")
	}
	defer conn.Close()

	client := miaow.NewChatServiceClient(conn)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < msgCount; i++ {
		message := new(miaow.Message)
		message.Body = fmt.Sprintf("%#08x", rand.Uint32())

		log.Printf("Sending: %s", message.Body)

		reply, err := client.SayHello(context.Background(), message)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Received: %s", reply.Body)
	}

	log.Print("Streaming...")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	stream, err := client.StreamHello(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		stream.Send(&miaow.Message{Body: "Abracadabras"})
		stream.Send(&miaow.Message{Body: "Wooohoo0ohoo"})
		stream.Send(&miaow.Message{Body: "Incorporated"})
		stream.CloseSend()
	}()

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}
			log.Printf("Received: %s", msg.Body)
		}
	}()

	<-ctx.Done()
}
