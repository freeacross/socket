package client

import (
	"context"
	"github.com/freeacross/socket"
	log "github.com/sirupsen/logrus"
	"net"
	"reflect"
)

func NewClient(t ...ClientParameterTemplate) *Client {
	var c = new(Client)
	c.Conf(t...)

	if c.ctx == nil {
		c.ctx = context.Background()
	}
	return c
}

type Client struct {
	_socket          socket.ISocket
	_listen          net.Conn
	isStop           bool
	ctx              context.Context
	name             string
	timeout          int
	network, address string
	Routers          [][2]interface{}
}

//func (s *Client) Route(rule interface{}, controller Controller) {
//	c._socket.Route(rule, controller)
//
//}

// Route 路由注册
func (c *Client) Route(rule interface{}, controller socket.Controller) {
	if reflect.TypeOf(rule).Implements(socket.FILTERITERN) {
		var arr [2]interface{}
		arr[0] = rule
		arr[1] = controller
		c.Routers = append(c.Routers, arr)
	} else {
		switch rule.(type) {
		default:
			// TODO 增加更人性化的路由选择
			log.Error("Something is wrong in Router:%!+v", rule)
		}
	}
}

func (c *Client) Run() {
	if c.ctx == nil {
		c.ctx = context.Background()
	}
	var (
		err error
	)

	c._listen, err = net.Dial(c.network, c.address)
	if err != nil {
		panic(err)
	}
	for {
		if c.isStop {
			break
		}
		log.Info("Waiting for clients")

		log.Infoln(c._listen.RemoteAddr().String(), " tcp connect success")
		// 如果此链接超过60秒没有发送新的数据，将被关闭
		c._socket = socket.NewSocket(c.ctx, "Client", &socket.Conn{c._listen}, c.timeout, c.Routers)
		c._socket.HandleConnection()
	}
}

func (c *Client) WriteData(data []byte) (n int, err error) {
	return c._socket.WriteData(data)
}

func (c *Client) Close() error {
	log.Infof("%s-%s: ready to exist\n", c.name, c._listen.RemoteAddr().String())

	c.isStop = true
	return c._listen.Close()
}

type ClientParameterTemplate func(s *Client)

func (c *Client) Conf(t ...ClientParameterTemplate) *Client {
	for _, f := range t {
		f(c)
	}
	return c
}

func Ctx(ctx context.Context) ClientParameterTemplate {
	return func(c *Client) {
		c.ctx = ctx
	}
}

//func ReceiveChannel(r chan []byte) ClientParameterTemplate {
//	return func(s *Client) {
//		c.ReceiveChannel = r
//	}
//}

func Name(n string) ClientParameterTemplate {
	return func(c *Client) {
		c.name = n
	}
}

func Timeout(t int) ClientParameterTemplate {
	return func(c *Client) {
		c.timeout = t
	}
}

func NetworkAddress(network, address string) ClientParameterTemplate {
	return func(c *Client) {
		c.network = network
		c.address = address
	}
}
