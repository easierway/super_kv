package super_kv

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

func CreateDelDataRequest(key []byte, value []byte) []byte {
	cmd := Command{
		Op: OP_DEL,
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
