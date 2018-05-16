package super_kv

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LeveDB_Engine struct {
	DataPath string
	Engine   *leveldb.DB
}

func CreateLevelDBEngine(dataPath string) (KV_Engine, error) {
	engine := LeveDB_Engine{}
	engine.DataPath = dataPath
	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		return nil, err
	}
	engine.Engine = db
	return &engine, nil
}

func (engine *LeveDB_Engine) Set(key []byte, value []byte) error {
	return engine.Engine.Put(key, value, nil)
}

func (engine *LeveDB_Engine) Get(key []byte) ([]byte, error) {
	return engine.Engine.Get(key, nil)
}

func (engine *LeveDB_Engine) Delete(key []byte) error {
	return engine.Engine.Delete(key, nil)
}
