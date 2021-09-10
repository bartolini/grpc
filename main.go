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
	log.Printf("Server received [%d]: %s", s.msgCount, msg.Body)
	reply := new(miaow.Message)
	reply.Body = fmt.Sprintf("=== PONG %s ===", msg.Body)
	return reply, nil
}

func (s *serverImplementation) StreamHello(stream miaow.ChatService_StreamHelloServer) error {
	log.Print("New stream...")
	defer log.Print("Goodbye my stream...")
	for {
		if msg, err := stream.Recv(); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		} else {
			s.msgCount++
			log.Printf("Stream received [%d]: %s", s.msgCount, msg.Body)
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
		log.Fatal(err)
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

	ctx, cancel := context.WithCancel(context.Background())

	stream, err := client.StreamHello(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for i := 0; i < msgCount; i++ {
			message := new(miaow.Message)
			message.Body = fmt.Sprintf("%#08x", rand.Uint32())
			log.Printf("Streaming: %s", message.Body)
			if err := stream.Send(message); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Millisecond)
		}
		stream.CloseSend()
		cancel()
	}()

	go func() {
		for {
			if msg, err := stream.Recv(); err != nil {
				return
			} else {
				log.Printf("Received: %s", msg.Body)
			}
		}
	}()

	<-ctx.Done()
}
