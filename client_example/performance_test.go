package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"super_kv"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/easierway/super_kv"
)

var conn net.Conn

const PAYLOAD_LEN = 2

func createConn() net.Conn {
	var err error
	addr := "127.0.0.1:9003"
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return conn
}

func ReceiveData1(conn net.Conn) ([]byte, error) {

	var respHeader [PAYLOAD_LEN]byte
	//conn.SetReadDeadline(time.Now().Add(READ_TIME_OUT))
	n, err := conn.Read([]byte(respHeader[:]))
	//var respHeader [2]byte

	//conn.SetReadDeadline(time.Now().Add(READ_TIME_OUT))
	//n, err := conn.Read([]byte(respHeader[:]))
	if err != nil {
		fmt.Println("Receive data error ", err)
		if err == io.EOF {
			fmt.Println("Receive data error is EOF")
			return nil, err
		}
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			return nil, err
		}
	}
	fmt.Println("Received1 header:", n, respHeader)
	if n == 0 {
		return nil, nil
	}
	if n < PAYLOAD_LEN {
		fmt.Println("less size data ", n, respHeader)
		return nil, err
	}

	dataLen := *(*uint16)(unsafe.Pointer(&respHeader))
	fmt.Println(dataLen)
	if dataLen > 0 {
		data := make([]byte, dataLen)
		n, err := io.ReadFull(conn, data)
		if n < int(dataLen) && err == nil {
			err = err
		}
		return data, err
	}
	return nil, err
}

func ReceiveResponse(conn net.Conn) {
	//super_kv.ReceiveData(conn)
	resp, err := super_kv.ReceiveServerResponse(conn)
	if err != nil {
		fmt.Println("response err", err)
		return
	}
	fmt.Println("response", resp.Ack, string(resp.Data))

}

func operationTest(wg *sync.WaitGroup, conn net.Conn) {

	start := time.Now()
	value := []byte("hello world 拨号操作，需要指定协议")

	for i := 0; i < 1; i++ {
		key := []byte("key" + strconv.Itoa(i))
		conn.Write(super_kv.CreateSetDataRequest(key, value))
		fmt.Println("Waiting for set response")
		ReceiveResponse(conn)
		conn.Write(super_kv.CreateGetDataRequest(key))
		fmt.Println("Waiting for get response")
		ReceiveResponse(conn)
		conn.Write(super_kv.CreateDelDataRequest(key, value))
		fmt.Println("Waiting for del response")
		ReceiveResponse(conn)
		conn.Write(super_kv.CreateGetDataRequest(key))
		fmt.Println("Waiting for get response")
		ReceiveResponse(conn)

	}
	conn.Close()
	end := time.Now()
	fmt.Println("time spent (ms):", end.Sub(start).Seconds()*1000)
	wg.Done()
}

func TestGet(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go operationTest(&wg, createConn())
	}
	wg.Wait()
}
