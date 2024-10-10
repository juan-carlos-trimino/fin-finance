package concurrency

/***
The actual Go implementation of channels integrates with the runtime scheduler to improve
performance. The channel's source code is located at
https://github.com/golang/go/blob/master/src/runtime/chan.go.
***/

import (
  "container/list"
  "sync"
)

/*
**
This implementation of a channel uses two semaphores, a mutex, and a linked list.
The actual Go implementation of channels integrates with the runtime scheduler to improve
performance; the source code can be found in Go's GitHub project under the runtime package located
at https://github.com/golang/go/blob/master/src/runtime/chan.go.
**
*/
type Channel[M any] struct { //Using generics.
  //Capacity semaphore to block sender when the buffer is full.
  senderSema *Semaphore
  //Buffer size semaphore to block receiver when the buffer is empty.
  receiverSema *Semaphore
  //Mutex to protect the shared list data structure.
  mutex sync.Mutex
  //Linked list use as a queue data structure (first in, first out [FIFO]).
  buffer *list.List
}

func NewChannel[M any](capacity int) *Channel[M] {
  return &Channel[M]{
    senderSema: NewSemaphore(capacity),
    receiverSema: NewSemaphore(0),
    mutex: sync.Mutex{},
    buffer: list.New(), //Empty linked list.
  }
}

/***
The Send() function needs to fulfill these three requirements:
(1) Block the goroutine if the buffer is full.
  If the senderSema is not full, the goroutine reduces the count by 1 and continues. Otherwise,
  it blocks.
  
  If when creating the channel an initial capacity of 0 (i.e., no buffer) is specified, the
  sender will block if a receiver is not present. This gives the same synchronous functionality
  of the default channel in Go; i.e., Go's channels are synchronous by default, meaning that the
  sender goroutine will block until there is a receiver goroutine ready to consume the message.
(2) Otherwise, safely add the message to the buffer.
  The mutex protects the queue from concurrent updates.
(3) If any receiver goroutines are blocked waiting for messages, resume one of them.
  The sender goroutine increments the count on the receiverSema; if there are goroutines waiting
  for messages, one will be resumed.
***/
func (c *Channel[M]) Send(message M) {
  c.senderSema.Acquire() //(1)
  c.mutex.Lock() //(2)
  c.buffer.PushBack(message)
  c.mutex.Unlock()
  c.receiverSema.Release() //(3)
}

/***
The Receive() function needs to satisfy the following requirements:
(1) Unblock a sender waiting because the senderSema is full.
  The receiver increments the count of the senderSema by 1. This will unblock a sender waiting on
  the senderSema.

  The reason for releasing the senderSema first is that the implementation needs to work when
  there is a zero-buffer channel.
(2) If the buffer is empty, block the receiver.
  The receiver tries to acquire the receiverSema. If the buffer is empty, the receiver goroutine
  will block.
(3) Otherwise, safely consume the next message from the buffer.
  Once the receiverSema unblocks the receiver, the goroutine reads and removes the next message
  from the buffer. The mutex protects the queue from concurrent updates.
***/
func (c *Channel[M]) Receive() M {
  c.senderSema.Release() //(1)
  c.receiverSema.Acquire() //(2)
  c.mutex.Lock()
  m := c.buffer.Remove(c.buffer.Front()).(M)
  c.mutex.Unlock()
  return m
}
