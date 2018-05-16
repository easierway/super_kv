package super_kv

const (
	ACK_OK                    byte = 0
	ACK_CONN_WAITING_TIME_OUT      = 1
	ACK_UNPACK_DATA_ERROR          = 2
	ACK_TOO_LARGE_RETURN_DATA      = 3
	ACK_FAILED                     = 4
	ACK_NO_SUCK_OPERATION          = 5
	ACK_FAILED_TO_RECEIVE          = 6

	OP_SET byte = 0
	OP_GET byte = 1
	OP_DEL byte = 2

	PAYLOAD_LEN int = 2 // lenOfRepsonseData (2 bytes - uint16)
)
