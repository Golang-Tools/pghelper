package proxy

import (
	log "github.com/Golang-Tools/loggerhelper"
	"github.com/go-pg/pg/v10"
)

//Callback redis操作的回调函数
type Callback func(cli *pg.DB) error

//pgProxy redis客户端的代理
type pgProxy struct {
	*pg.DB
	parallelcallback bool
	callBacks        []Callback
}

// New 创建一个新的数据库客户端代理
func New() *pgProxy {
	proxy := new(pgProxy)
	return proxy
}

// IsOk 检查代理是否已经可用
func (proxy *pgProxy) IsOk() bool {
	if proxy.DB == nil {
		return false
	}
	return true
}

//SetConnect 设置连接的客户端
//@params cli UniversalClient 满足redis.UniversalClient接口的对象的指针
func (proxy *pgProxy) SetConnect(cli *pg.DB) error {
	if proxy.IsOk() {
		return ErrProxyAllreadySettedClient
	}

	proxy.DB = cli
	if proxy.parallelcallback {
		for _, cb := range proxy.callBacks {
			go func(cb Callback) {
				err := cb(proxy.DB)
				if err != nil {
					log.Error("regist callback get error", log.Dict{"err": err})
				} else {
					log.Debug("regist callback done")
				}
			}(cb)
		}
	} else {
		for _, cb := range proxy.callBacks {
			err := cb(proxy.DB)
			if err != nil {
				log.Error("regist callback get error", log.Dict{"err": err})
			} else {
				log.Debug("regist callback done")
			}
		}
	}
	return nil
}

//InitFromOptions 从配置条件初始化代理对象
func (proxy *pgProxy) InitFromOptions(options *pg.Options) error {
	cli := pg.Connect(options)
	return proxy.SetConnect(cli)
}

//InitFromOptionsParallelCallback 从配置条件初始化代理对象,并行执行回调函数
func (proxy *pgProxy) InitFromOptionsParallelCallback(options *pg.Options) error {
	cli := pg.Connect(options)
	proxy.parallelcallback = true
	return proxy.SetConnect(cli)
}

//InitFromURL 从URL条件初始化代理对象
func (proxy *pgProxy) InitFromURL(url string) error {
	options, err := pg.ParseURL(url)
	if err != nil {
		return err
	}
	return proxy.InitFromOptions(options)
}

//InitFromURLParallelCallback 从URL条件初始化代理对象
func (proxy *pgProxy) InitFromURLParallelCallback(url string) error {
	options, err := pg.ParseURL(url)
	if err != nil {
		return err
	}
	return proxy.InitFromOptionsParallelCallback(options)
}

// Regist 注册回调函数,在init执行后执行回调函数
//如果对象已经设置了被代理客户端则无法再注册回调函数
func (proxy *pgProxy) Regist(cb Callback) error {
	if proxy.IsOk() {
		return ErrProxyAllreadySettedClient
	}
	proxy.callBacks = append(proxy.callBacks, cb)
	return nil
}

//Proxy 默认的redis代理对象
var Proxy = New()
