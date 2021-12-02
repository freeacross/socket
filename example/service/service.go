package main

import (
	"context"
	"github.com/freeacross/socket/server"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
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

func main() {
	s := server.NewServer(
		server.Name("client"),
		server.Ctx(context.Background()),
		server.NetworkAddress("tcp", "localhost:6060"),
		server.Timeout(-1),
	)

	var controller Controller
	kvs := make(map[string]string)
	kvs["msgType"] = "send SMS"
	s.Route(kvs, &controller)

	s.Run()
}
