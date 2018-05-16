package main

import (
	"super_kv"
	"time"
)

func main() {
	server := super_kv.Server{
		DataPath:           "path/to/db",
		Port:               9003,
		ConnBufSize:        1000,
		NumOfConnHandler:   100,
		ConnWaitingTimeout: time.Millisecond * 10,
		LogConfig:          "seelog.xml",
	}
	server.StartServer()
}
