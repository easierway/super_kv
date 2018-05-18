package main

import (
	"fmt"

	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	super_kv "github.com/easierway/super_kv"
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

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func operationTest(wg *sync.WaitGroup, conn net.Conn, t *testing.T) {

	start := time.Now()
	value := []byte("hello world 拨号操作，需要指定协议")
	clientChan := super_kv.CreateChannel(conn)
	var (
		resp *super_kv.Response
		err  error
	)
	for i := 0; i < 1; i++ {
		key := []byte("key" + strconv.Itoa(i))
		resp, err = clientChan.Set(key, value)
		checkError(err, t)
		fmt.Println("Set Response ", resp)
		resp, err = clientChan.Get(key)
		checkError(err, t)
		fmt.Println("Get", resp)
		expectedStr := string(value)
		actualStr := string(resp.Data)
		fmt.Println("String Content:", actualStr)
		if actualStr != expectedStr {
			t.Errorf("Expected value is %s, but actual value is %s\n", expectedStr, actualStr)
		}
		resp, err = clientChan.Delete(key)
		checkError(err, t)
		fmt.Println("Delete Response ", resp)
		resp, err = clientChan.Get(key)
		checkError(err, t)
		fmt.Println("Get", resp)
		if len(resp.Data) != 0 {
			t.Errorf("Expected value is %d, but actual value is %d\n",
				0, len(resp.Data))
		}
	}
	conn.Close()
	end := time.Now()
	fmt.Println("time spent (ms):", end.Sub(start).Seconds()*1000)
	wg.Done()
}

var server super_kv.Server

func StartServer() {
	server = super_kv.Server{
		DataPath:           "path/to/db",
		Port:               9003,
		ConnBufSize:        1000,
		NumOfConnHandler:   100,
		ConnWaitingTimeout: time.Millisecond * 10,
		LogConfig:          "seelog.xml",
	}
	server.StartServer()
}

func TestBasicOperation(t *testing.T) {
	var wg sync.WaitGroup
	go StartServer()
	time.Sleep(time.Second * 2)
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go operationTest(&wg, createConn(), t)
	}
	wg.Wait()
	server.StopServer()
}
