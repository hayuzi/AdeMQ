package packet

import (
	"encoding/binary"
	"errors"
	"net"
)

const (
	ConstHeadSize      = 4
	ConstBufferFullErr = "消息长度超出缓冲区上限"
)

type ReadableEventChan chan []byte
type WritableEventChan chan []byte

type TcpConn struct {
	Conn          net.Conn
	readBuf       []byte
	readStart     int
	readEnd       int
	maxReadBufLen int
	Closed        bool
	ReadableEventChan
	WritableEventChan
}

// New 封装已经建立的TCP连接
// 	@params coon net.Conn 表示建立的连接
// 	@params readBufLen int 表示读取缓冲区的长度
// 	@params maxReadBufLen int 表示读取缓冲区的上限长度。注意，单条消息超出该上限，我们会设计触发异常，并返回
func New(conn net.Conn, readBufLen int, maxReadBufLen int) *TcpConn {
	// 接收缓冲区默认最大为10M
	if maxReadBufLen == 0 {
		maxReadBufLen = 1024 * 1024 * 10
	}
	readBuf := make([]byte, readBufLen)
	readChan := make(chan []byte, 10)
	writeChan := make(chan []byte, 10)
	return &TcpConn{
		conn,
		readBuf,
		0,
		0,
		maxReadBufLen,
		false,
		readChan,
		writeChan,
	}
}

func (tc *TcpConn) Close() {
	tc.Closed = true
	defer close(tc.WritableEventChan)
	defer close(tc.ReadableEventChan)
}

// ReadFromConn 从conn里面读取数据，conn可能阻塞
func (tc *TcpConn) ReadFromConn() (int, error) {
	tc.readBufLeftShift()
	// 在缓冲区不能承载整个消息体的时候，我们需要对缓冲区扩容, 或者直接抛出异常
	if tc.readEnd >= len(tc.readBuf) {
		// 如果超出缓冲区的最大上限，我们还是应该做限制
		if tc.readEnd >= tc.maxReadBufLen {
			return 0, errors.New(ConstBufferFullErr)
		}
		newBuf := make([]byte, len(tc.readBuf)*2)
		copy(newBuf, tc.readBuf)
		tc.readBuf = newBuf
	}
	// 此处如果传入读取的缓冲区空闲长度为0，会陷入死循环， 所以前面做了扩容以及异常处理
	n, err := tc.Conn.Read(tc.readBuf[tc.readEnd:])
	if err != nil {
		return n, err
	}
	tc.readEnd += n
	return n, nil
}

// leftShift 将读取缓冲区的有用字节前移
func (tc *TcpConn) readBufLeftShift() {
	if tc.readStart == 0 {
		return
	}
	// copy 用于将内容从一个数组切片复制到另一个数组切片。如果加入的两个数组切片不一样大，就会按其中较小的那个数组切片的元素个数进行复制。
	// 注意 copy是浅拷贝
	copy(tc.readBuf, tc.readBuf[tc.readStart:tc.readEnd])
	tc.readEnd -= tc.readStart
	tc.readStart = 0
}

// ReadBufLen 获取当前读取缓冲区的内容长度
func (tc *TcpConn) ReadBufLen() int {
	return tc.readEnd - tc.readStart
}

// Seek 返回n个字节，而不产生移位
func (tc *TcpConn) Seek(n int) ([]byte, error) {
	if tc.readEnd-tc.readStart >= n {
		buf := tc.readBuf[tc.readStart : tc.readStart+n]
		return buf, nil
	}
	return nil, errors.New("not enough")
}

// Read 舍弃offset个字段，读取n个字段
func (tc *TcpConn) Read(offset, n int) []byte {
	tc.readStart += offset
	buf := tc.readBuf[tc.readStart : tc.readStart+n]
	tc.readStart += n
	return buf
}

// SendMessageToChan 向连接发送消息
func (tc *TcpConn) SendMessageToChan(content []byte) error {
	if tc.Closed {
		return errors.New("连接发送通道已经关闭")
	}
	tc.WritableEventChan <- content
	return nil
}

// SendMessageDirect 向连接发送消息
func (tc *TcpConn) SendMessageDirect(content []byte) error {
	headBytes := make([]byte, ConstHeadSize)
	contentSize := len(content)
	headBytes = tc.IntToBytes(contentSize)
	_, err := tc.Conn.Write(append(headBytes, content...))
	if err != nil {
		return err
	}
	return nil
}

// IntToBytes 整形转换成字节
func (tc *TcpConn) IntToBytes(n int) []byte {
	var b = make([]byte, ConstHeadSize)
	// (32位下如果超出整型上限可能有错误)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

// BytesToInt 字节转换成整型
func (tc *TcpConn) BytesToInt(b []byte) int {
	// (32位下如果超出整型上限可能有错误)
	return int(binary.BigEndian.Uint32(b))
}
