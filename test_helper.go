package super_kv

import (
	"testing"
)

func checkTestError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

//For checking basic engine operations (set/get/delete)
func checkBasicOps(engine KV_Engine, t *testing.T) {
	key := []byte("Key")
	value := []byte("Hello")
	err := engine.Set(key, value)
	checkTestError(err, t)
	v, err := engine.Get(key)
	checkTestError(err, t)
	if string(v) != string(value) {
		t.Errorf("The expected value is %s, but the value is %s",
			string(value), string(v))
	}
	err = engine.Delete(key)
	checkTestError(err, t)
	v1, err := engine.Get(key)
	if v1 != nil && len(v1) != 0 {
		t.Errorf("The expected value is nil/[], but the value is %v",
			v1)
	}
}
