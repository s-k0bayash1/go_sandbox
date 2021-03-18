package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/s-k0bayash1/go_sandbox/grpc_stream/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"
)

func main() {
	listenPort, err := net.Listen("tcp", ":50011")
	if err != nil {
		log.Fatalln(err)
	}
	server := grpc.NewServer()
	service := &chatService{}
	// 実行したい実処理をseverに登録する
	pb.RegisterChatServiceServer(server, service)
	log.Println("start server")
	server.Serve(listenPort)
}

type chatService struct {
	connectionsMutex sync.RWMutex
	connections []*connection
	messagesMutex sync.RWMutex
	messages []string
}

type connection struct {
	id int64
	stream pb.ChatService_ConnectServer
	err chan error
}

func (c *connection) Send(msg string) {
	if err := c.stream.Send(&pb.Message{Msg: msg}); err != nil {
		c.err <- err
	}
}

func (c *chatService) Connect(req *pb.Connection, stream pb.ChatService_ConnectServer) error {
	conn := func() *connection {
		c.connectionsMutex.Lock()
		defer c.connectionsMutex.Unlock()
		for _, connection := range c.connections {
			if connection.id == req.Id {
				return nil
			}
		}
		conn := &connection{
			id:     req.Id,
			stream: stream,
			err: make(chan error),
		}
		c.connections = append(c.connections, conn)
		return conn
	}()
	if conn == nil {
		return status.Error(codes.AlreadyExists, fmt.Sprintf("already exists id: %d", req.Id))
	}
	log.Println(fmt.Sprintf("add connection id=%d", req.Id))
	c.messagesMutex.RLock()
	for _, message := range c.messages {
		conn.Send(message)
	}
	c.messagesMutex.RUnlock()
	return <- conn.err
}

func (c *chatService) SendMessage(_ context.Context, req *pb.Message) (*empty.Empty, error) {
	var wg sync.WaitGroup
	c.messagesMutex.Lock()
	defer c.messagesMutex.Unlock()
	c.messages = append(c.messages, req.Msg)
	for _, conn := range c.connections {
		wg.Add(1)
		conn := conn
		go func() {
			conn.Send(req.Msg)
			wg.Done()
			return
		}()
	}
	wg.Wait()
	return &empty.Empty{}, nil
}

func (c *chatService) Close(_ context.Context, req *pb.Connection) (*empty.Empty, error) {
	c.connectionsMutex.Lock()
	for i, connection := range c.connections {
		if connection.id == req.Id {
			connection.err <- status.Error(codes.Canceled, "connection closed")
			c.connections = append(c.connections[:i], c.connections[i+1:]...)
			break
		}
	}
	c.connectionsMutex.Unlock()
	log.Println(fmt.Printf("close connection id=%d", req.Id))
	return &empty.Empty{}, nil
}
