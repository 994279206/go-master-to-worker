package util

import uuid "github.com/satori/go.uuid"

var HttpPort string
var Ip string
var GrpcPort string
var NodeName string
var HeartBeatTime int

func HasUuid() string {
	str := uuid.NewV4().String()
	return str

}
