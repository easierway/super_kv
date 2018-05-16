package super_kv

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	READ_TIME_OUT  time.Duration = time.Second * 10
	WRITE_TIME_OUT time.Duration = time.Second * 1
)

var NoSuchOperationError error = errors.New("No such operation")

type ConnHandler struct {
	kvEngine KV_Engine
	connBuf  <-chan net.Conn
}

func convertErrorToAck(err error) byte {
	if err != nil {
		return ACK_FAILED
	}
	return ACK_OK
}

func (handler *ConnHandler) processCommand(cmd *Command) (resp []byte, ack byte) {
	switch cmd.Op {
	case OP_SET:
		err := handler.kvEngine.Set([]byte(cmd.Params[0]), []byte(cmd.Params[1]))
		return nil, convertErrorToAck(err)
	case OP_GET:
		fmt.Println("GET KEY", string(cmd.Params[0]))
		data, err := handler.kvEngine.Get([]byte(cmd.Params[0]))
		if ack := convertErrorToAck(err); ack != ACK_OK {
			return nil, ack
		} else {
			fmt.Println("GET DATA", string(data))
			return data, ack
		}
	case OP_DEL:
		err := handler.kvEngine.Delete([]byte(cmd.Params[0]))
		return nil, convertErrorToAck(err)
	default:
		return nil, ACK_NO_SUCK_OPERATION
	}
}

func (handler *ConnHandler) acceptConn() {
	for {
		conn := <-handler.connBuf
		handler.handleWithoutPanic(conn)
	}
}

func (handler *ConnHandler) handleWithoutPanic(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			if logger != nil {
				logger.Error("recovered from error", err)
			} else {
				fmt.Println("recovered from error", err)
			}
		}
	}()
	handler.handle(conn)
}

func (handler *ConnHandler) send(conn net.Conn, data []byte) error {
	//conn.SetDeadline(time.Now().Add(WRITE_TIME_OUT))
	n, err := conn.Write(data)
	fmt.Println("Send", n, data)
	if err != nil {
		fmt.Println("Failed to send", err)
		logger.Error("Failed to send data", err)
	}
	return err
}

func (handler *ConnHandler) handlerDataReceivingError(conn net.Conn, err error, ack byte) {
	output := packResponse(ack, []byte(err.Error()))
	handler.send(conn, output)
}

func (handler *ConnHandler) handle(conn net.Conn) {
	defer func(conn net.Conn) {
		if conn != nil {
			fmt.Println("Connection closed")
			conn.Close()
		}
	}(conn)
	cnt := 0
	for {
		cnt++
		fmt.Println("Received cnt", cnt)
		data, err := ReceiveData(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("Failed to receive data", err)
			handler.handlerDataReceivingError(conn, err, ACK_FAILED_TO_RECEIVE)
			break
		}
		if data == nil {
			continue
		}
		cmd, err := UnpackData(data)
		if err != nil {
			logger.Error("Failed to unpack data", err)
			handler.handlerDataReceivingError(conn, err, ACK_UNPACK_DATA_ERROR)
			break
		}
		resp, ack := handler.processCommand(cmd)
		output := packResponse(ack, resp)
		sendErr := handler.send(conn, output)
		if sendErr != nil {
			logger.Error("Failed to send data ", err)
			break
		}
	}
}
