package linear

import (
	"errors"
)

// RingQueue 环形队列
type RingQueue struct {
	qElem   []interface{} // 初始化的动态分配存储空间
	maxSize int           // 最大队列长度
	front   int           // 头指针
	rear    int           // 尾指针
}

// NewRingQueue 创建一个循环队列
func NewRingQueue(maxSize int) *RingQueue {
	return &RingQueue{
		qElem:   make([]interface{}, maxSize),
		maxSize: maxSize,
		front:   0,
		rear:    0,
	}
}

// Length 获取队列长度
func (q *RingQueue) Length() int {
	return (q.rear - q.front + q.maxSize) % q.maxSize
}

// EnQueue 插入元素
func (q *RingQueue) EnQueue(e interface{}) bool {
	// 少用一个元素空间便于判定头尾以及是否队列满
	// 队列满返回失败
	if (q.rear+1)%q.maxSize == q.front {
		return false
	}
	q.qElem[q.rear] = e
	q.rear = (q.rear + 1) % q.maxSize
	return true
}

// DeQueue 从队列中弹出元素
func (q *RingQueue) DeQueue() (interface{}, error) {
	if q.front == q.rear {
		return nil, errors.New("ring queue is empty")
	}
	e := q.qElem[q.front]
	q.front = (q.front + 1) % q.maxSize
	return e, nil
}

// IsEmpty 是否为空
func (q *RingQueue) IsEmpty() bool {
	return q.front == q.rear
}

// IsFull 是否为满
func (q *RingQueue) IsFull() bool {
	return (q.rear+1)%q.maxSize == q.front
}

// FetchAllElem 获取所有内容
func (q *RingQueue) FetchAllElem() []interface{} {
	n := q.Length()
	data := make([]interface{}, n)
	for i := 0; i < n; i++ {
		idx := (q.front + i) % q.maxSize
		data[i] = q.qElem[idx]
	}
	return data
}
