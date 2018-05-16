package super_kv

import (
	"testing"
)

func checkError(t *testing.T, msg string, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestPackUnpack(t *testing.T) {
	params := [][]byte{
		[]byte("Hello World"),
		[]byte("Super KV"),
	}
	cmd := Command{
		Op:     1,
		Params: params,
	}
	data, err := PackData(&cmd)
	checkError(t, "error happened when packing data", err)
	cmd1, err := UnpackData(data)
	checkError(t, "error happened when unpacking data", err)
	if cmd1.Op != cmd.Op {
		t.Errorf("expected cmd.Op is %d, but the value is %d\n", cmd.Op, cmd1.Op)
	}
	if string(cmd1.Params[0]) != string(cmd.Params[0]) {
		t.Errorf("expected cmd.Parmas[0] is %s, but the value is %s\n",
			string(cmd.Params[0]), string(cmd1.Params[0]))
	}
	if string(cmd1.Params[1]) != string(cmd.Params[1]) {
		t.Errorf("expected cmd.Parmas[1] is %s, but the value is %s\n",
			string(cmd.Params[1]), string(cmd1.Params[1]))
	}
}
