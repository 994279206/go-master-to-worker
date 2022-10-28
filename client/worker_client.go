package client

import (
	"context"
	"fmt"
	"go_grpc/core"
	"go_grpc/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Worker struct {
	conn *grpc.ClientConn
	c    core.NodeServiceClient
}

var worker *Worker

func (w *Worker) Init() error {
	var err error
	clientAddress := fmt.Sprintf("%v:%v", util.Ip, util.GrpcPort)
	log.Println(util.NodeName, "grpc address", clientAddress)
	w.conn, err = grpc.Dial(clientAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}
	w.c = core.NewNodeServiceClient(w.conn)
	return nil
}
func (w *Worker) Start() {
	log.Println("start worker node")
	ctx := context.Background()
	go w.HeartBeat()
	stream, _ := w.c.AssignTask(ctx, &core.Request{})
	for {
		res, err := stream.Recv()
		if err != nil {
			return
		}
		log.Println("received command: ", res.Data)

		parts := strings.Split(res.Data, " ")
		if err = exec.CommandContext(ctx, parts[0], parts[1:]...).Run(); err != nil {
			log.Panic(err)
		}
	}
}
func (w *Worker) HeartBeat() {
	uuid := util.HasUuid()
	ctx := context.Background()
	for {
		_, err := w.c.ReportStatus(ctx, &core.Request{Action: uuid})
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Duration(util.HeartBeatTime) * time.Second)
	}
}

func NewWorkNode() *Worker {
	if worker == nil {
		worker = &Worker{}
		if err := worker.Init(); err != nil {
			log.Panic(err)
		}
	}
	return worker
}
