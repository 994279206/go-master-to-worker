package service

import (
	"context"
	"go_grpc/core"
)

type NodeServer struct {
	core.UnimplementedNodeServiceServer
	Channel chan string
}

func (n *NodeServer) ReportStatus(ctx context.Context, request *core.Request) (*core.Response, error) {
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
		}
	}
	return server
}
