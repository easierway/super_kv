package super_kv

/*
 Install gorocksdb
CGO_CFLAGS="-I/usr/local/Cellar/rocksdb/5.12.4/include" \
CGO_LDFLAGS="-L/usr/local/Cellar/rocksdb/5.12.4 -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy" \
  go get github.com/tecbot/gorocksdb
*/
import (
	"github.com/tecbot/gorocksdb"
)

type RocksDB_Engine struct {
	DataPath string
	Engine   *gorocksdb.DB
	ReadOpt  *gorocksdb.ReadOptions
	WriteOpt *gorocksdb.WriteOptions
}

func CreateRocksDBEngine(dataPath string) (KV_Engine, error) {
	ro := gorocksdb.NewDefaultReadOptions()
	wo := gorocksdb.NewDefaultWriteOptions()
	//bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	//	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	//opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, "path/to/rocksdb")
	if err != nil {
		return nil, err
	}

	return &RocksDB_Engine{
		DataPath: dataPath,
		Engine:   db,
		ReadOpt:  ro,
		WriteOpt: wo,
	}, nil
}

func (engine *RocksDB_Engine) Set(key []byte, value []byte) error {
	return engine.Engine.Put(engine.WriteOpt, key, value)
}

func (engine *RocksDB_Engine) Get(key []byte) ([]byte, error) {
	value, err := engine.Engine.Get(engine.ReadOpt, key)
	defer value.Free()
	return value.Data(), err
}

func (engine *RocksDB_Engine) Delete(key []byte) error {
	return engine.Engine.Delete(engine.WriteOpt, key)
}
