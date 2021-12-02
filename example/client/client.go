package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/freeacross/socket"
)

func senderMsg(conn socket.Conn) {
	msg := "fork message"
	conn.WriteData([]byte(msg))
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil{
		log.Printf("error:%v", err)
	}
	log.Printf("%s receive data string:%+v \n", conn.RemoteAddr().String(), string(buffer[:n]))
}

func main() {
	// server := "localhost:6060"
	// tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	// 	os.Exit(1)
	// }

	conn, err := socket.Dial("tcp", ":6060")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	for i := 1; i < 100; i++ {
		senderMsg(conn)
		time.Sleep(time.Second)
	}

}
