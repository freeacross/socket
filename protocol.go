//通讯协议处理，主要处理封包和解包的过程
package socket

import (
	"bytes"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

// TODO modify
const (
	constHeader       = "Char(2)"
	constHeaderLength = len(constHeader)

	constTail       = "Char(3)"
	constTailLength = len(constTail)
)

//封包
func packet(message []byte) []byte {
	return append(append([]byte(constHeader), message...), []byte(constTail)...)
}

//解包
func unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i++ {
		if length <= i+constHeaderLength {
			break
		}
		if string(buffer[i:i+constHeaderLength]) == constHeader {
			// 进入正文
			if index := bytes.Index(buffer[i+constHeaderLength:], []byte(constTail)); index == -1 {
				break
			} else {
				data := buffer[i+constHeaderLength : i+constHeaderLength+index]
				readerChannel <- data
				log.Debugln("------>", string(data))
				i += constHeaderLength + index
				//back step
				i--
			}
		}
	}

	if i == length {
		return make([]byte, 0)
	}

	return buffer[i:]
}

// 暂时不需要
//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
