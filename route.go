package socket

import (
	"io"
	"log"
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
		if pred.(IdentifyRequestF)(data) /*判断是否是对应注册的handle*/{
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
		Log("origin:[%s];conn.WriteError()", string(data), err)
	}
}

func reader(conn Conn, readerChannel <-chan []byte, timeout int) {
	for {
		select {
		case data := <-readerChannel:
			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			// handle request
			handleRequest(conn, data)
			break
		case <-time.After(time.Duration(timeout) * time.Second):
			conn.Close()
			Log("connection is closed.")
			return
		}
	}
}

// HandleConnection 处理长连接
func HandleConnection(conn Conn, timeout int) {
	//声明一个临时缓冲区，用来存储被截断的数据
	var tmpBuffer []byte

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(conn, readerChannel, timeout)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				continue
			}
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				Log("exit goroutine.")
				return
			}
			Log(conn.RemoteAddr().String(), " connection error: ", err, reflect.TypeOf(err))
			return
		}
		tmpBuffer = unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
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
		Log("Something is wrong in Router")
	}
}

func init() {
	Routers = make([][2]interface{}, 0, 10)
}
