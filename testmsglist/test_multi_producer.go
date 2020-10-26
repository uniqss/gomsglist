package main

import (
	"fmt"
	"github.com/uniqss/gomsglist"
	"strings"
)

const TEST_PRODUCER_CONSUMER_COUNT = 10000
const TEST_MSG_COUNT_PER_PRODUCER = 100000

var producerDone [TEST_PRODUCER_CONSUMER_COUNT]bool
var producerDoneAll = false
var consumedmsgs [TEST_PRODUCER_CONSUMER_COUNT][TEST_MSG_COUNT_PER_PRODUCER]bool
var ML = gomsglist.NewSafeMsgList()

func producer(idx int64) {
	var msg int64
	var i int64 = 0
	for ; i < TEST_MSG_COUNT_PER_PRODUCER; i++ {
		msg = (idx << 32) | i
		ML.Put(msg)
	}
	producerDone[idx] = true
	for i := 0; i < TEST_PRODUCER_CONSUMER_COUNT; i++ {
		if !producerDone[i] {
			return
		}
	}

	producerDoneAll = true
}

func consumer(idx int) {
	for !producerDoneAll {
		var err error = gomsglist.ErrNoNode
		for err != nil {
			var msg interface{}
			msg, err = ML.Pop()
			if err == nil {
				msgint64 := msg.(int64)
				consumerIdx := int32(msgint64 >> 32)
				msgint := int32(msgint64 & 0xffffffff)
				if consumedmsgs[consumerIdx][msgint] {
					panic("this should not happen")
				}
				consumedmsgs[consumerIdx][msgint] = true
			}
		}
	}
}

func check() bool {
	for idx := 0; idx < TEST_PRODUCER_CONSUMER_COUNT; idx++ {
		for i := 0; i < TEST_MSG_COUNT_PER_PRODUCER; i++ {
			if !consumedmsgs[idx][i] {
				return false
			}
		}
	}

	return true
}

func main() {
	var i int64 = 0
	for ; i < TEST_PRODUCER_CONSUMER_COUNT; i++ {
		go producer(i)
	}

	var i2 = 0
	//for ; i2 < TEST_PRODUCER_CONSUMER_COUNT; i2++ {
	go consumer(i2)
	//}

	var input string
	for {
		fmt.Scanln(&input)
		input = strings.ToLower(input)
		if input == "e" || input == "exit" {
			break
		}
		if input == "c" || input == "check" {
			fmt.Println("check:", check())
		}
	}
}
