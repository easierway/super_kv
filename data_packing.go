package super_kv

import (
	"errors"
	"fmt"
	"io"
	"net"
	//	"time"
	"unsafe"
)

var (
	OperationTimeoutErr          error = errors.New("Operation time out")
	ResponseSizeIsLessErr        error = errors.New("Response size is less 3")
	ReceivedDataIsNotEnoughErr   error = errors.New("Received data is not enough")
	DeclaredPayloadSizeIsZeroErr error = errors.New("Declared data size is 0")
)

type Command struct {
	Op     byte
	Params [][]byte
}

type Response struct {
	Ack  byte
	Data []byte
}

func PackRequest(data []byte) []byte {
	lenOfData := len(data)
	lenOfDataByte := *(*[2]byte)(unsafe.Pointer(&lenOfData))
	request := make([]byte, lenOfData+2)
	request[0] = lenOfDataByte[0]
	request[1] = lenOfDataByte[1]
	copy(request[2:], data)
	return request
}

/**
  |Len(2)|ACK(1)|data(Len)|
**/
func packResponse(ack byte, data []byte) []byte {
	var response []byte
	if data == nil {
		response = make([]byte, 3)
		response[0] = 1 //only ACK in package
		response[1] = 0
		response[2] = ack
		return response
	} else {
		lenOfData := len(data) + 1
		if lenOfData > 65536 {
			return packResponse(ACK_TOO_LARGE_RETURN_DATA, nil)
		}
		response = make([]byte, 2+lenOfData)
		lenBytes := *(*[2]byte)(unsafe.Pointer(&lenOfData))
		response[0] = lenBytes[0]
		response[1] = lenBytes[1]
		response[2] = ack
		copy(response[3:], data)
		fmt.Println(response)
		return response
	}
}

func ReceiveData(conn net.Conn) ([]byte, error) {
	var respHeader [PAYLOAD_LEN]byte
	//conn.SetReadDeadline(time.Now().Add(READ_TIME_OUT))
	n, err := conn.Read([]byte(respHeader[:]))
	if err != nil {
		fmt.Println("Receive data error ", err)
		if err == io.EOF {
			fmt.Println("Receive data error is EOF")
			return nil, err
		}
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			return nil, OperationTimeoutErr
		}
	}
	fmt.Println("Received header:", n, respHeader)
	if n == 0 {
		return nil, nil
	}
	if n < PAYLOAD_LEN {
		fmt.Println("less size data ", n, respHeader)
		return nil, ResponseSizeIsLessErr
	}

	dataLen := *(*uint16)(unsafe.Pointer(&respHeader))
	fmt.Println(dataLen)
	if dataLen > 0 {
		data := make([]byte, dataLen)
		n, err := conn.Read(data)
		if n < int(dataLen) && err == nil {
			err = ReceivedDataIsNotEnoughErr
		}
		return data, err
	}

	return nil, DeclaredPayloadSizeIsZeroErr
}

func ReceiveServerResponse(conn net.Conn) (*Response, error) {
	payload, err := ReceiveData(conn)
	if err != nil {
		return nil, err
	}
	ack := payload[0]
	responseData := Response{}
	responseData.Ack = ack
	data := make([]byte, len(payload)-1)
	fmt.Println("payload", payload[1:])
	copy(data, payload[1:])
	responseData.Data = data
	return &responseData, nil
}

func UnpackData(rawData []byte) (*Command, error) {
	cmd := Command{}
	cmd.Op = rawData[0]

	offset := 1
	var numParamsByte [1]byte
	numParamsByte[0] = rawData[offset]
	numParam := *(*uint8)(unsafe.Pointer(&numParamsByte))
	cmd.Params = make([][]byte, numParam)
	var lenBytes [2]byte

	for i := 0; i < int(numParam); i++ {
		offset++
		lenBytes[0] = rawData[offset]
		offset++
		lenBytes[1] = rawData[offset]
		dataLen := *(*uint16)(unsafe.Pointer(&lenBytes))
		offset++
		param := ([]byte)(rawData[offset : offset+int(dataLen)])
		offset += int(dataLen) - 1
		cmd.Params[i] = param
	}

	return &cmd, nil
}

func PackData(cmd *Command) ([]byte, error) {
	lenOfOutput := 1 /*op*/ + 1 /*numOfParam*/
	for i := 0; i < len(cmd.Params); i++ {
		lenOfOutput += len(cmd.Params[i]) + 2
	}
	output := make([]byte, lenOfOutput)
	output[0] = cmd.Op
	output[1] = (byte)(len(cmd.Params))
	var (
		lenParamBytes [2]byte
		lenParam      uint16
	)
	offset := 1
	for i := 0; i < len(cmd.Params); i++ {
		lenParam = (uint16)(len(cmd.Params[i]))
		lenParamBytes = *(*[2]byte)(unsafe.Pointer(&lenParam))
		offset++
		output[offset] = lenParamBytes[0]
		offset++
		output[offset] = lenParamBytes[1]
		for j := 0; j < len(cmd.Params[i]); j++ {
			offset++
			output[offset] = cmd.Params[i][j]
		}
	}
	return output, nil

}
