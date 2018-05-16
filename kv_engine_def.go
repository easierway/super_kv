package super_kv

type KV_Engine interface {
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
}
