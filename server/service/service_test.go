package service

import (
	"encoding/binary"
	"fmt"
	"github.com/AdeMQ/protocol/packet"
	"net"
	"testing"
	"time"
)

func TestConnectAndMsgSend(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:10601")
	if err != nil {
		fmt.Println("连接失败", err)
	}
	defer conn.Close()

	var headSize int
	var headBytes = make([]byte, packet.ConstHeadSize)
	s := "hello world！hello world！hello world！hello world！hello world！hello world！hello world！hello world！" +
		"hello world！hello world！hello world！hello world！"
	content := []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint32(headBytes, uint32(headSize))

	t.Log(headBytes)
	t.Log(string(content))
	conn.Write(headBytes)
	conn.Write(content)

	s = "hello go"
	content = []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint32(headBytes, uint32(headSize))
	conn.Write(headBytes)
	conn.Write(content)

	s = "hello service"
	content = []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint32(headBytes, uint32(headSize))
	conn.Write(headBytes)
	conn.Write(content)

	time.Sleep(time.Second * 5)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Log(err.Error())
	}
	t.Log(string(buf[packet.ConstHeadSize:n]))

}
