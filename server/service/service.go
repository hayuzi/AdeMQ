package service

import (
	"fmt"
	"github.com/AdeMQ/protocol/packet"
	"log"
	"net"
	"time"
)

type Config struct {
	Address   string `yaml:"address" json:"address"`
	BufLen    int    `yaml:"bufLen" json:"bufLen"`
	BufMaxLen int    `yaml:"bufMaxLen" json:"bufMaxLen"`
}

// Run 启动服务
func Run(conf *Config) (err error) {

	// 开启TCP的端口监听
	ln, err := net.Listen("tcp", conf.Address)
	if err != nil {
		log.Println("Error start listen", err.Error())
		return
	}
	for {
		// 等待客户端建立连接
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accept connect", err.Error())
			continue
		}
		// 开启新的协程处理连接
		go handleConnection(conn, conf)
	}

}

// 连接处理函数
func handleConnection(conn net.Conn, conf *Config) {
	defer closeConnection(conn)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))

	// TCP数据包边界问题（俗称TCP粘包问题）
	// 由于 TCP 本身是面向字节流的，无法理解上层的业务数据，所以在底层是无法保证数据包不被拆分和重组的，
	// 这个问题只能通过上层的应用协议栈设计来解决，根据业界的主流协议的解决方案，一般有三种：
	// 		消息定长、设置消息边界、将消息分为消息头和消息体

	// 这里我们将采用消息头+消息体的方法来 确定消息边界
	// TODO 边界界定，获取完整消息之后触发完整的 receive 事件，相应的，之后给客户端回传消息也将采用该方式处理
	// 客户端回收结果也会是类似的处理方式
	var (
		// TcpConn 实现了 io.reader 接口，我们可以用自己封装的 buffer 来处理
		tcpConn     = packet.New(conn, conf.BufLen, conf.BufMaxLen)
		headBuf     []byte
		contentSize int
		contentBuf  []byte
	)

	defer tcpConn.Close()

	// 开启向该连接发送消息的协程, 阻塞监听消息, 如果连接关闭，则退出
	go handleWriteConnection(tcpConn)

	// 循环阻塞读取消息, 读取到的消息追加存储到消息体中, 待消息收满之后, 发送给程序处理
	for {
		_, err := tcpConn.ReadFromConn()
		if err != nil {
			log.Println("Error reading", err.Error())
			if err.Error() == packet.ConstBufferFullErr {
				// 因为需要立即返回，此处就直接发送到连接中
				_ = tcpConn.SendMessageDirect([]byte(packet.ConstBufferFullErr))
			}
			return
		}
		for {
			// 刚开始的消息默认是消息头, 消息头一般设计为占用2字节或者4字节的长度, 用来保存整个消息的长度
			// 消息长度在设计的时候可以单纯的为消息体长度，可以为 包含消息头以及消息体的总长度
			// 此处设计采用消息长度 为 单纯的消息体长度
			headBuf, err = tcpConn.Seek(packet.ConstHeadSize)
			if err != nil {
				break
			}
			// 提取包头中存储的包体长度: 二进制大端字节序列无符号整数转换为int
			contentSize = tcpConn.BytesToInt(headBuf)
			// 如果缓冲区中的内容长度超过或者等于 消息头+消息体长度，那么后面相当于读取到了消息体的消息
			if tcpConn.ReadBufLen() >= contentSize+packet.ConstHeadSize {
				// 将完整的消息体内容读取到缓冲区，进行后续处理
				// TODO 分发数据并处理，处理过程中涉及到的消息返回需要设计整体架构
				contentBuf = tcpConn.Read(packet.ConstHeadSize, contentSize)
				fmt.Println(string(contentBuf))
				continue
			}
			break
		}
	}
}

// 关闭连接
func closeConnection(conn net.Conn) {
	_ = conn.Close()
}

// 连接消息发送处理函数
func handleWriteConnection(tcpConn *packet.TcpConn) {
	select {
	case msg := <-tcpConn.WritableEventChan:
		if msg == nil {
			// chan关闭了
			log.Println("Error Writing 连接已经关闭")
			return
		}
		if tcpConn.Closed {
			goto End
		}
		if err := tcpConn.SendMessageDirect(msg); err != nil {
			log.Println("Error Writing 消息发送失败", err.Error())
		}
	}
End:
}
