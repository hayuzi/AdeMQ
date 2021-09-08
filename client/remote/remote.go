package remote

import (
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"time"
)

const (
	ConstHeadSize      = 4
	ConstBufferFullErr = "消息长度超出缓冲区上限"
)

type Remote struct {
	closed             bool        // 链接是否关闭
	Conn               net.Conn    // rP连接
	maxReadBufLen      int         // 最大接受缓冲区长度
	readBuf            []byte      // 缓冲区
	readStart          int         // 缓冲区数据开始位置
	readEnd            int         // 缓冲区数据结束位置
	RequestChan        chan []byte // 远程请求发送通道
	RequestChanClosed  bool        // 远程请求队列是否关闭
	ResponseChan       chan []byte // 远程结果通道（ 命令处理函数发送之后阻塞接收 ）
	ResponseChanClosed bool        // 远程结果队列是否关闭
}

var address = flag.String("address", "127.0.0.1:10601", "远程服务端地址")

func NewRemote() *Remote {
	// TODO 远程连接的地址，需要通过，启动命令的时候给
	conn, err := net.DialTimeout("tcp", *address, 5*time.Second)
	if err != nil {
		panic("服务器连接失败")
	}
	return &Remote{
		closed:             false,
		Conn:               conn,
		maxReadBufLen:      1024 * 1024 * 10,
		readBuf:            make([]byte, 1024*16),
		readStart:          0,
		readEnd:            0,
		RequestChan:        make(chan []byte),
		RequestChanClosed:  false,
		ResponseChan:       make(chan []byte),
		ResponseChanClosed: false,
	}
}

func (r *Remote) Init() {
	go r.HandleConnWrite()
	go r.HandleConnRead()
	go r.HandleHeartBeat()
}

func (r *Remote) HandleConnRead() {
	// 循环阻塞读取消息, 读取到的消息追加存储到消息体中, 待消息收满之后, 发送给程序处理
	var (
		headBuf     []byte
		contentSize int
		contentBuf  []byte
	)
	for {
		_, err := r.ReadFromConn()
		if err != nil {
			log.Println("Error reading", err.Error())
			// 数据超出限制
			if err.Error() == ConstBufferFullErr {
				r.Close()
				panic(ConstBufferFullErr)
			}
			return
		}
		for {
			// 刚开始的消息默认是消息头, 消息头一般设计为占用2字节或者4字节的长度, 用来保存整个消息的长度
			// 消息长度在设计的时候可以单纯的为消息体长度，可以为 包含消息头以及消息体的总长度
			// 此处设计采用消息长度 为 单纯的消息体长度
			headBuf, err = r.Seek(ConstHeadSize)
			if err != nil {
				break
			}
			// 提取包头中存储的包体长度: 二进制大端字节序列无符号整数转换为int
			contentSize = r.BytesToInt(headBuf)
			// 如果缓冲区中的内容长度超过或者等于 消息头+消息体长度，那么后面相当于读取到了消息体的消息
			if r.ReadBufLen() >= contentSize+ConstHeadSize {
				// 将完整的消息体内容读取到缓冲区, 并发送到结果通道
				contentBuf = r.Read(ConstHeadSize, contentSize)
				if r.RequestChanClosed {
					break
				}
				r.ResponseChan <- contentBuf
			}
			break
		}
	}
}

func (r *Remote) HandleConnWrite() {
	// 请求通道有数据就写入到远程
	for req := range r.RequestChan {
		if req != nil {
			// 发送数据
			err := r.sendMsgDirect(req)
			if err == nil {
				continue
			}
			if !r.ResponseChanClosed {
				r.ResponseChan <- []byte(err.Error())
			} else {
				log.Println("remote cmd send failed")
			}
		}
	}
}

func (r *Remote) HandleHeartBeat() {
	err := r.sendMsgDirect([]byte("heart beat"))
	if err != nil {
		log.Println("heart beat send err ", err.Error())
	}
}

// leftShift 将读取缓冲区的有用字节前移
func (r *Remote) readBufLeftShift() {
	if r.readStart == 0 {
		return
	}
	// copy 用于将内容从一个数组切片复制到另一个数组切片。如果加入的两个数组切片不一样大，就会按其中较小的那个数组切片的元素个数进行复制。
	// 注意 copy是浅拷贝
	copy(r.readBuf, r.readBuf[r.readStart:r.readEnd])
	r.readEnd -= r.readStart
	r.readStart = 0
}

// ReadBufLen 获取当前读取缓冲区的内容长度
func (r *Remote) ReadBufLen() int {
	return r.readEnd - r.readStart
}

// Seek 返回n个字节，而不产生移位
func (r *Remote) Seek(n int) ([]byte, error) {
	if r.readEnd-r.readStart >= n {
		buf := r.readBuf[r.readStart : r.readStart+n]
		return buf, nil
	}
	return nil, errors.New("not enough")
}

// Read 舍弃offset个字段，读取n个字段
func (r *Remote) Read(offset, n int) []byte {
	r.readStart += offset
	buf := r.readBuf[r.readStart : r.readStart+n]
	r.readStart += n
	return buf
}

// SendMsgToRequestChan 向连接发送消息
func (r *Remote) SendMsgToRequestChan(content []byte) error {
	if r.RequestChanClosed {
		return errors.New("remote request chan closed")
	}
	r.RequestChan <- content
	return nil
}

// GetResponseFromChan 向连接发送消息
func (r *Remote) GetResponseFromChan() ([]byte, error) {
	if r.ResponseChanClosed {
		return nil, errors.New("remote response chan closed")
	}
	select {
	case msg := <-r.ResponseChan:
		return msg, nil
	case <-time.After(0 * time.Second):
		return []byte(""), nil
	}
}

// sendMessageDirect 向连接发送消息
func (r *Remote) sendMsgDirect(content []byte) error {
	headBytes := make([]byte, ConstHeadSize)
	contentSize := len(content)
	headBytes = r.IntToBytes(contentSize)
	_, err := r.Conn.Write(append(headBytes, content...))
	if err != nil {
		return err
	}
	return nil
}

// IntToBytes 整形转换成字节
func (r *Remote) IntToBytes(n int) []byte {
	var b = make([]byte, ConstHeadSize)
	// (32位下如果超出整型上限可能有错误)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

// BytesToInt 字节转换成整型
func (r *Remote) BytesToInt(b []byte) int {
	// (32位下如果超出整型上限可能有错误)
	return int(binary.BigEndian.Uint32(b))
}

func (r *Remote) Close() {
	r.closed = true
	r.RequestChanClosed = true
	r.ResponseChanClosed = true
	defer r.closeConn()
	defer close(r.RequestChan)
	defer close(r.ResponseChan)
}

func (r *Remote) closeConn() {
	_ = r.Conn.Close()
}

// ReadFromConn 从conn里面读取数据，conn可能阻塞
func (r *Remote) ReadFromConn() (int, error) {
	r.readBufLeftShift()
	// 在缓冲区不能承载整个消息体的时候，我们需要对缓冲区扩容, 或者直接抛出异常
	if r.readEnd >= len(r.readBuf) {
		// 如果超出缓冲区的最大上限，我们还是应该做限制
		if r.readEnd >= r.maxReadBufLen {
			return 0, errors.New(ConstBufferFullErr)
		}
		newBuf := make([]byte, len(r.readBuf)*2)
		copy(newBuf, r.readBuf)
		r.readBuf = newBuf
	}
	// 此处如果传入读取的缓冲区空闲长度为0，会陷入死循环， 所以前面做了扩容以及异常处理
	n, err := r.Conn.Read(r.readBuf[r.readEnd:])
	if err != nil {
		if err == io.EOF {
			r.Close()
		}
		return n, err
	}
	r.readEnd += n
	return n, nil
}
