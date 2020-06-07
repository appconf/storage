package storage

import (
	"fmt"
	"sync"

	"github.com/appconf/log"
)

//Data 配置存储器与Engine之间的交流数据
type Data struct {
	Key   string
	Value interface{}
	Err   error
}

//Driver 自定义Storage需要实现此接口进行注册
type Driver interface {
	Open(log.Logger, map[string]interface{}) (Storage, error)
}

//Storage 配置存储器的接口
type Storage interface {
	//Get 获取指定keys的信息
	Get(keys []string) (ch chan []Data, err error)
	//Stop 停止Storage
	Stop() error
	//Success 当key所在模板更新成功时使用此接口进行通知
	Success(template string, kvs []map[string]interface{})
	//Error 当key所在模板渲染失败时使用此接口进行通知storage
	Error(template string, kvs []map[string]interface{}, err error)
}

var (
	storages     = make(map[string]Driver)
	storagesLock sync.RWMutex
)

//Register 注册storage名称, storage的构造函数
func Register(name string, driver Driver) {
	storagesLock.Lock()
	defer storagesLock.Unlock()

	if _, dup := storages[name]; dup {
		panic("storage: register called twice for storage " + name)
	}

	storages[name] = driver
}

//List 已注册的存储器列表
func List() (s []string) {
	storagesLock.Lock()
	defer storagesLock.Unlock()

	for k := range storages {
		s = append(s, k)
	}
	return
}

//New 根据name从已注册的storages获取对应的构造函数，创建storage
func New(name string, logger log.Logger, cfg map[string]interface{}) (Storage, error) {
	storagesLock.Lock()
	defer storagesLock.Unlock()

	driver, ok := storages[name]
	if !ok {
		return nil, fmt.Errorf("storage: not found storage %s", name)
	}

	return driver.Open(logger, cfg)
}
