package socket

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"reflect"
	"time"
)

type IdentifyRequestF func(requestData []byte) bool

// handle request
func handleRequest(conn Conn, data []byte) {

	for _, v := range Routers {
		pred := v[0]
		act := v[1]
		if pred.(IdentifyRequestF)(data) /*判断是否是对应注册的handle*/ {
			log.Debugf("handleRequest:%+v", pred)
			// yes, to handle this request.
			result := act.(Controller).Handle(data)
			_, err := writeResult(conn, result)
			if err != nil {
				Log("conn.WriteResult()", err)
			}
			return
		}
	}

	_, err := writeError(conn, "1111", "不能处理此类型的业务")
	if err != nil {
		log.Errorf("origin:[%s];conn.WriteError:%v", string(data), err)
	}
}

func (s *Socket) click() <-chan time.Time {
	// block ever ever
	if s.timeout == -1 {
		return nil
	}
	return time.After(time.Duration(s.timeout) * time.Second)
}

func (s *Socket) setDeadline() {
	if s.timeout == -1 {
		return
	}
	s.conn.SetDeadline(time.Now().Add(time.Duration(s.timeout) * time.Second))
}

func (s *Socket) handle() {
	for {
		select {
		case data := <-s.ReceiveChannel:
			s.Log("receive one package that was unpacked")
			s.setDeadline()
			// handle request
			handleRequest(s.conn, data)
			break
		case <-s.click():
			s.conn.Close()
			Log("connection is closed.")
			return
		}
	}
}

func NewSocket(name string, conn Conn, timeout int) *Socket {
	var socket = Socket{
		conn:           conn,
		timeout:        timeout,
		Name:           name,
		ReceiveChannel: make(chan []byte, 16),
	}
	return &socket
}

type Socket struct {
	conn           Conn
	timeout        int // -1:表示不超时
	Name           string
	ReceiveChannel chan []byte
}

func (s *Socket) WriteData(data []byte) (n int, err error) {
	return s.conn.WriteData(data)
}

func (s *Socket) Log(v ...interface{}) {
	log.Println(s.Name, ":", fmt.Sprintln(v...))
}

// HandleConnection 处理长连接
func (s *Socket) HandleConnection() {
	//声明一个临时缓冲区，用来存储被截断的数据
	var tmpBuffer []byte

	//声明一个管道用于接收解包的数据
	go s.handle()

	buffer := make([]byte, 1024)

	// 开始循环读数据
	for {
		n, err := s.conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				//s.Log("reading EOF")
				continue
			}
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				Log("exit goroutine.")
				return
			}
			s.Log(s.conn.RemoteAddr().String(), " connection error: ", err, reflect.TypeOf(err))
			return
		}
		// 处理粘包
		tmpBuffer = unpack(append(tmpBuffer, buffer[:n]...), s.ReceiveChannel)
	}
}

func Log(v ...interface{}) {
	log.Println(v...)
}

// Routers 路由
// 二维数组:第二个为对应的处理程序；第一个为判断此请求是否是需要处理的
var Routers [][2]interface{}

// Route 路由注册
func Route(rule interface{}, controller Controller) {
	switch rule.(type) {
	case IdentifyRequestF:
		{
			var arr [2]interface{}
			arr[0] = rule
			arr[1] = controller
			Routers = append(Routers, arr)
		}
		// TODO 增加更人性化的路由选择
	default:
		Log("Something is wrong in Router:%!+v", rule)
	}
}

func init() {
	Routers = make([][2]interface{}, 0, 10)
}
