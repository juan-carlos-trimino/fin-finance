package concurrency

import (
  "sync"
)

/***
Read-preferred readers-writer lock.
As long as there are readers, writers will be blocked; i.e., reader goroutines are STARVING the
writer goroutines by not allowing them access to the shared resource.
https://go.dev/ref/mem
***/
type RPReadWriteLock struct {
  //Count the number of reader goroutines currently in the critical section.
  readersCounter int
  //Mutex for synchronizing readers access.
  readersLock sync.Mutex
  //Mutex for blocking any writers access.
  writersLock sync.Mutex
}

func (rw *RPReadWriteLock) ReadLock() {
  //Synchronizes access so that only one goroutine is allowed at any time.
  rw.readersLock.Lock()
  rw.readersCounter++
  //If first reader, attempt to lock writers.
  if rw.readersCounter == 1 {
    rw.writersLock.Lock()
  }
  rw.readersLock.Unlock()
}

func (rw *RPReadWriteLock) ReadUnlock() {
  rw.readersLock.Lock()
  rw.readersCounter--
  //If last reader, unlock writers.
  if rw.readersCounter == 0 {
    rw.writersLock.Unlock()
  }
  rw.readersLock.Unlock()
}

func (rw *RPReadWriteLock) WriteLock() {
  //Any writer access requires a lock on writersLock.
  rw.writersLock.Lock()
}

func (rw *RPReadWriteLock) WriteUnlock() {
  //Release the lock.
  rw.writersLock.Unlock()
}
