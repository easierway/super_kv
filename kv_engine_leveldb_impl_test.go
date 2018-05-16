package super_kv

import (
	"testing"
)

func TestBasicOp(t *testing.T) {
	engine, err := CreateLevelDBEngine("path/to/db")
	checkTestError(err, t)
	checkBasicOps(engine, t)
}
