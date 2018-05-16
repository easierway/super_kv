package super_kv

import (
	"testing"
)

func TestBasicOp(t *testing.T) {
	engine, err := CreateRocksDBEngine("path/to/db1")
	checkTestError(err, t)
	checkBasicOps(engine, t)
}
