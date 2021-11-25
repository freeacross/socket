package socket

import (
	"log"
	"testing"
)

var inputString = "Char(2)hello server, I'm clientChar(3)"

func TestUnpack(t *testing.T) {
	var readerChannel = make(chan []byte, 16000)

	ret := unpack([]byte(inputString), readerChannel)

	log.Println(string(ret))
}
