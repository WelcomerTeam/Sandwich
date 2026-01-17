package internal

import "sync"

type ChannelBuffer[T any] struct {
	Out    chan T
	buffer []T
	cond   *sync.Cond
}

func NewChannelBuffer[T any]() *ChannelBuffer[T] {
	channelBuffer := &ChannelBuffer[T]{
		Out:    make(chan T),
		buffer: make([]T, 0),
		cond:   sync.NewCond(&sync.Mutex{}),
	}

	go channelBuffer.run()

	return channelBuffer
}

func (cb *ChannelBuffer[T]) Push(item T) {
	cb.cond.L.Lock()
	cb.buffer = append(cb.buffer, item)
	cb.cond.Signal() // wake a waiter
	cb.cond.L.Unlock()
}

func (cb *ChannelBuffer[T]) Len() int {
	cb.cond.L.Lock()
	length := len(cb.buffer)
	cb.cond.L.Unlock()

	return length
}

func (cb *ChannelBuffer[T]) run() {
	for {
		cb.cond.L.Lock()
		for len(cb.buffer) == 0 {
			cb.cond.Wait() // sleep until Push signals
		}

		item := cb.buffer[0]
		cb.buffer = cb.buffer[1:]
		cb.cond.L.Unlock()

		cb.Out <- item
	}
}
