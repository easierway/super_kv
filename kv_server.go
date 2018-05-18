package super_kv

import (
	"net"
	"strconv"
	"time"

	"github.com/cihub/seelog"
)

const DATA_PATH = "path/to/db"

var logger seelog.LoggerInterface

type Server struct {
	DataPath           string
	Port               int
	ConnBufSize        int
	connBuf            chan net.Conn
	NumOfConnHandler   int
	ConnWaitingTimeout time.Duration
	LogConfig          string
	stopChan           chan struct{}
}

func (server *Server) startHandlerPool() {
	kvEngine, err := CreateLevelDBEngine(server.DataPath)
	checkFatalError(err)
	server.connBuf = make(chan net.Conn, server.ConnBufSize)
	for i := 0; i < server.NumOfConnHandler; i++ {
		handler := ConnHandler{}
		handler.kvEngine = kvEngine
		handler.connBuf = server.connBuf
		go handler.acceptConn()
	}
}

func (server *Server) StartServer() error {
	var err error
	logger, err = seelog.LoggerFromConfigAsFile("seelog.xml")
	server.stopChan = make(chan struct{})
	checkFatalError(err)
	logger.Info("Starting server")
	server.startHandlerPool()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(server.Port))
	checkFatalError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkFatalError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Failed to accept connection ", err)
			continue
		}
		select {
		case <-server.stopChan:
			break
		case server.connBuf <- conn:
		case <-time.After(server.ConnWaitingTimeout):
			data := packResponse(ACK_CONN_WAITING_TIME_OUT, nil)
			conn.Write(data)
			conn.Close()
		}

	}
}

func (server *Server) StopServer() {
	close(server.stopChan)
}

func checkFatalError(err error) {
	if err != nil {
		logger.Error("Failed to start server ", err)
		logger.Flush()
		panic(err)
	}
}
