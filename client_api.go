package super_kv

import (
	"net"
)

type Channel struct {
	conn net.Conn
}

func CreateSetDataRequest(key []byte, value []byte) []byte {
	cmd := Command{
		Op: OP_SET,
		Params: [][]byte{
			key,
			value,
		},
	}
	data, err := PackData(&cmd)
	if err != nil {
		panic(err)
	}

	return PackRequest(data)
}

func CreateDelDataRequest(key []byte) []byte {
	cmd := Command{
		Op: OP_DEL,
		Params: [][]byte{
			key,
		},
	}
	data, err := PackData(&cmd)
	if err != nil {
		panic(err)
	}
	return PackRequest(data)
}

func CreateGetDataRequest(key []byte) []byte {
	cmd := Command{
		Op: OP_GET,
		Params: [][]byte{
			key,
		},
	}
	data, err := PackData(&cmd)
	if err != nil {
		panic(err)
	}

	return PackRequest(data)
}

func (c *Channel) Set(key []byte, value []byte) (*Response, error) {
	c.conn.Write(CreateSetDataRequest(key, value))
	return ReceiveServerResponse(c.conn)
}

func (c *Channel) Get(key []byte) (*Response, error) {
	c.conn.Write(CreateGetDataRequest(key))
	return ReceiveServerResponse(c.conn)
}

func (c *Channel) Delete(key []byte) (*Response, error) {
	c.conn.Write(CreateDelDataRequest(key))
	return ReceiveServerResponse(c.conn)
}

func CreateChannel(conn net.Conn) *Channel {
	return &Channel{conn}
}
