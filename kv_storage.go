package super_kv

type KV_Storage interface {
	Set(key interface{}, value interface{})
}
