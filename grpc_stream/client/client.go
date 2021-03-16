package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	pb "go_sandbox/grpc_stream/chat"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"os/signal"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 50011), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	id := flag.Int64("id", 1, "id")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := pb.NewChatServiceClient(conn)
	defer c.Close(context.Background(), &pb.Connection{Id: *id})
	done := make(chan interface{})
	go func() {
		if err := Connect(ctx, c, *id); err != nil {
			log.Println(err)
			cancel()
		}
		done <- struct {}{}
		return
	}()
	go func() {
		if err := SendMessage(ctx, c); err != nil {
			log.Println(err)
			cancel()
		}
		done <- struct {}{}
		return
	}()
	quit := make(chan os.Signal)
	defer close(quit)

	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		return
	case <-done:
		return
	}
}

func Connect(ctx context.Context, c pb.ChatServiceClient, id int64) (err error) {
	stream, err := c.Connect(ctx, &pb.Connection{Id: id})
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}
		log.Println(resp.Msg)
	}
}

func SendMessage(ctx context.Context, c pb.ChatServiceClient) (err error) {
	stdin := bufio.NewScanner(os.Stdin)
	for {
		stdin.Scan()
		msg := stdin.Text()
		_, err = c.SendMessage(ctx, &pb.Message{Msg: msg})
		if err != nil {
			return
		}
		if msg == "exit" {
			break
		}
	}
	return
}
