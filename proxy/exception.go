package proxy

import (
	"errors"
)

//ErrProxyAllreadySettedClient 代理已经设置过redis客户端对象
var ErrProxyAllreadySettedClient = errors.New("代理不能重复设置客户端对象")

//ErrProxyNotYetSettedClient 代理还未设置客户端对象
// var ErrProxyNotYetSettedClient = errors.New("代理还未设置客户端对象")
