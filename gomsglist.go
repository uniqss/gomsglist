package gomsglist

import (
	"errors"
	"runtime"
	"sync/atomic"
	"unsafe"
)

type _SafeListNode struct {
	next  unsafe.Pointer
	value interface{}
}

func newNode(data interface{}) unsafe.Pointer {
	return unsafe.Pointer(&_SafeListNode{
		nil,
		data,
	})
}

type SafeMsgList struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewSafeMsgList() *SafeMsgList {
	node := unsafe.Pointer(newNode(nil))
	return &SafeMsgList{
		node,
		node,
	}
}

func (sl *SafeMsgList) Put(data interface{}) {
	newNode := newNode(data)
	var tail unsafe.Pointer

	for {
		tail = sl.tail
		next := (*_SafeListNode)(tail).next

		if next != nil {
			atomic.CompareAndSwapPointer(&sl.tail, tail, next)
		} else {
			if atomic.CompareAndSwapPointer(&(*_SafeListNode)(sl.tail).next, nil, newNode) {
				break
			}
		}
		runtime.Gosched()
	}

	atomic.CompareAndSwapPointer(&sl.tail, tail, newNode)
}

var ErrNoNode = errors.New("no node")

func (sl *SafeMsgList) Pop() (interface{}, error) {
	for {
		head := sl.head
		tail := sl.tail

		next := (*_SafeListNode)(head).next

		if head == tail {
			if next == nil {
				return nil, ErrNoNode
			}
			atomic.CompareAndSwapPointer(&sl.tail, tail, next)
		} else {
			if atomic.CompareAndSwapPointer(&sl.head, head, next) {
				return (*_SafeListNode)(next).value, nil
			}
		}

		runtime.Gosched()
	}
}

func (sl *SafeMsgList) IsEmpty() bool {
	head := sl.head
	tail := sl.tail

	next := (*_SafeListNode)(head).next
	if head == tail {
		if next == nil {
			return true
		}
	}

	return false
}
