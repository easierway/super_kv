package super_kv

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB_Engine struct {
	DataPath string
	Engine   *leveldb.DB
}

func CreateLevelDBEngine(dataPath string) (KV_Engine, error) {
	engine := LevelDB_Engine{}
	engine.DataPath = dataPath
	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		return nil, err
	}
	engine.Engine = db
	return &engine, nil
}

func (engine *LevelDB_Engine) Set(key []byte, value []byte) error {
	return engine.Engine.Put(key, value, nil)
}

func (engine *LevelDB_Engine) Get(key []byte) ([]byte, error) {
	return engine.Engine.Get(key, nil)
}

func (engine *LevelDB_Engine) Delete(key []byte) error {
	return engine.Engine.Delete(key, nil)
}
