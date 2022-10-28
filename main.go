package main

import (
	flag "github.com/spf13/pflag"
	"go_grpc/client"
	"go_grpc/util"
	"log"
)

func init() {
	parse()
	initLog()

}

func parse() {
	flag.StringVarP(&util.NodeName, "node", "n", "", "This is which node")
	flag.StringVarP(&util.HttpPort, "HPort", "p", "9090", "http port")
	flag.StringVarP(&util.GrpcPort, "GPort", "g", "9091", "grpc port")
	flag.StringVarP(&util.Ip, "Ip", "i", "127.0.0.1", "server address ip")
	flag.IntVarP(&util.HeartBeatTime, "heartbeat", "", 5, "client heart beat time")
	flag.Parse()
}
func initLog() {
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("")
}
func main() {
	switch util.NodeName {
	case "master":
		client.NewMasterNode().Start()
	case "worker":
		client.NewWorkNode().Start()
	default:
		log.Panic("invalid node name")
	}
}
