package service

import (
	"context"
	"go_grpc/core"
	"time"
)

type NodeServer struct {
	core.UnimplementedNodeServiceServer
	Channel chan string
	Count   int
	Uid     map[string]int64
}

func (n *NodeServer) ReportStatus(ctx context.Context, request *core.Request) (*core.Response, error) {
	uuid := request.Uid
	if uuid != "" {
		if _, ok := n.Uid[uuid]; !ok {
			n.Count += 1
		}
		n.Uid[uuid] = time.Now().Unix()

	}
	return &core.Response{Data: "ok"}, nil

}
func (n *NodeServer) AssignTask(request *core.Request, server core.NodeService_AssignTaskServer) error {
	for {
		select {
		case c := <-n.Channel:
			if err := server.Send(&core.Response{Data: c}); err != nil {
				return err
			}
		}
	}
}

var server *NodeServer

func NewNodeServer() *NodeServer {
	if server == nil {
		server = &NodeServer{
			Channel: make(chan string),
			Uid:     make(map[string]int64),
		}
	}
	return server
}
