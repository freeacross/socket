package server

import (
	"context"
	"github.com/freeacross/socket"
	log "github.com/sirupsen/logrus"
	"net"
	"reflect"
)

func NewServer(t ...ServerParameterTemplate) *Server {
	var s = new(Server)
	s.Conf(t...)

	if s.ctx == nil {
		s.ctx = context.Background()
	}
	var (
		err error
	)
	s._listen, err = net.Listen(s.network, s.address)
	if err != nil {
		panic(err)
	}
	return s
}

type Server struct {
	_socket          socket.ISocket
	_listen          net.Listener
	isStop           bool
	ctx              context.Context
	name             string
	timeout          int
	network, address string
	Routers          [][2]interface{}
}

//func (s *Server) Route(rule interface{}, controller Controller) {
//	s._socket.Route(rule, controller)
//
//}

// Route 路由注册
func (s *Server) Route(rule interface{}, controller socket.Controller) {
	if reflect.TypeOf(rule).Implements(socket.FILTERITERN) {
		var arr [2]interface{}
		arr[0] = rule
		arr[1] = controller
		s.Routers = append(s.Routers, arr)
	} else {
		switch rule.(type) {
		default:
			// TODO 增加更人性化的路由选择
			log.Error("Something is wrong in Router:%!+v", rule)
		}
	}
}

func (s *Server) Run() {
	if s.ctx == nil {
		s.ctx = context.Background()
	}
	for {
		if s.isStop {
			break
		}
		log.Info("Waiting for clients")
		conn, err := s._listen.Accept()
		if err != nil {
			continue
		}
		log.Infoln(conn.RemoteAddr().String(), " tcp connect success")
		// 如果此链接超过60秒没有发送新的数据，将被关闭
		s._socket = socket.NewSocket(s.ctx, "Server", &socket.Conn{conn}, s.timeout, s.Routers)
		go s._socket.HandleConnection()
	}
}

func (s *Server) Close() error {
	log.Infof("%s-%s: ready to exist\n", s.name, s._listen.Addr().String())
	s.isStop = true
	return s._listen.Close()
}

type ServerParameterTemplate func(s *Server)

func (s *Server) Conf(t ...ServerParameterTemplate) *Server {
	for _, f := range t {
		f(s)
	}
	return s
}

func Ctx(ctx context.Context) ServerParameterTemplate {
	return func(s *Server) {
		s.ctx = ctx
	}
}

//func ReceiveChannel(r chan []byte) ServerParameterTemplate {
//	return func(s *Server) {
//		s.ReceiveChannel = r
//	}
//}

func Name(n string) ServerParameterTemplate {
	return func(s *Server) {
		s.name = n
	}
}

func Timeout(t int) ServerParameterTemplate {
	return func(s *Server) {
		s.timeout = t
	}
}

func NetworkAddress(network, address string) ServerParameterTemplate {
	return func(s *Server) {
		s.network = network
		s.address = address
	}
}
