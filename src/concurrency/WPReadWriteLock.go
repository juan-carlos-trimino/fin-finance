package concurrency

import (
  "sync"
)

type WPReadWriteLock struct {
  //Indicate if a writer goroutine is holding the write lock.
  writerActive bool
  //Count the number of reader goroutines currently holding the read lock.
  readersCounter int
  //The number of writer goroutines currently waiting.
  writersWaiting int
  cond *sync.Cond
}

func NewWPReadWriteLock() *WPReadWriteLock {
  return &WPReadWriteLock{
    writerActive: false,
    readersCounter: 0,
    writersWaiting: 0,
    cond: sync.NewCond(&sync.Mutex{}),
  }
}


