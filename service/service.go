package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"time"

	"github.com/freeacross/socket"
)

type Controller struct {
}

func (this *Controller) Handle(req []byte) interface{} {
	log.Printf("request:%s", string(req))
	if time.Now().Unix()%2 == 0 {
		return "another success"
	}
	return "success"
}

func CheckError(err error) {
	if err != nil {
		log.Printf("Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func init() {
	var controller Controller
	kvs := make(map[string]string)
	kvs["msgType"] = "send SMS"
	socket.Route(kvs, &controller)
}

func main() {
	netListen, err := net.Listen("tcp", "localhost:8080")
	CheckError(err)
	defer netListen.Close()
	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Log(conn.RemoteAddr().String(), " tcp connect success")
		// 如果此链接超过6秒没有发送新的数据，将被关闭
		go socket.NewSocket("server", socket.Conn{conn}, -1).HandleConnection()
	}
}
