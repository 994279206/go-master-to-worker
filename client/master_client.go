package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_grpc/core"
	"go_grpc/service"
	"go_grpc/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"net/http"
	"time"
)

var master *Master

var loadData struct {
	Cmd string `json:"cmd"`
}

type Master struct {
	api     *gin.Engine
	ln      *net.TCPListener
	svr     *grpc.Server
	nodeSvr *service.NodeServer
}

var kasp = keepalive.ServerParameters{
	Time:    5 * time.Second, // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout: 1 * time.Second, // Wait 1 second for the ping ack before assuming the connection is dead
}

func (m *Master) checkClient() {
	for {
		nowTime := time.Now().Unix()
		offline := 0
		for uid, timeStamp := range m.nodeSvr.Uid {
			if nowTime-timeStamp > int64(util.HeartBeatTime*2) {
				delete(m.nodeSvr.Uid, uid)
				m.nodeSvr.Count -= 1
				offline += 1
			}
		}
		log.Println(fmt.Sprintf("在线客户端:%v, 离线客户端数量:%v", m.nodeSvr.Count, offline))
		time.Sleep(time.Duration(util.HeartBeatTime*2) * time.Second)
	}

}
func (m *Master) Init() error {
	grpcAddress := fmt.Sprintf("%v:%v", util.Ip, util.GrpcPort)
	log.Println(util.NodeName, "grpc list on", grpcAddress)
	address, err := net.ResolveTCPAddr("tcp", grpcAddress)
	if err != nil {
		return err
	}
	m.ln, err = net.ListenTCP("tcp", address)
	if err != nil {
		return err
	}
	m.nodeSvr = service.NewNodeServer()
	m.svr = grpc.NewServer(grpc.KeepaliveParams(kasp))
	core.RegisterNodeServiceServer(m.svr, m.nodeSvr)
	m.api = gin.Default()
	m.api.POST("/tasks", m.taskHandler)
	return nil
}
func (m *Master) taskHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&loadData); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	m.nodeSvr.Channel <- loadData.Cmd
	c.AbortWithStatus(http.StatusOK)
}

func (m *Master) Start() {
	log.Println("start master node")
	go m.runGin()
	go m.checkClient()
	err := m.svr.Serve(m.ln)
	if err != nil {
		log.Panic(err)
	}
	m.svr.Stop()

}
func (m *Master) runGin() {
	httpAddress := fmt.Sprintf("%v:%v", util.Ip, util.HttpPort)
	log.Println(util.NodeName, "http list on", httpAddress)
	err := m.api.Run(httpAddress)
	if err != nil {
		log.Panic(err)
	}
}

func NewMasterNode() *Master {
	if master == nil {
		master = &Master{}
		if err := master.Init(); err != nil {
			log.Panic(err)
		}
	}
	return master
}
