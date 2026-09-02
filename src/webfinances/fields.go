package webfinances

//
//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "encoding/json"
  "errors"
  "finance/concurrency/LockManager"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gposu"
  "os"
  "path/filepath"
  "sync"
  "time"
)

var (
  mainDir string = "/finances"
  /***
  In Go, built-in maps are not thread-safe. If one goroutine is writing to currentFields while another goroutine is reading from
  currentFields, the Go application will crash instantly with a non-recoverable error: fatal error: concurrent map read and map write.
  ***/
  currentFields = make(map[string]*UserSession)
  //Allow multiple goroutines to read a shared resource simultaneously, but grants exclusive access to a single writer.
  currentFieldsLock sync.RWMutex
  UserLockMgr = lockmanager.NewLockManager()  //Global manager instance.
)

type fields struct {
  //Make the pointers unexported so that clients can't interact with them directly but only via exported methods.
  miscellaneous *miscellaneousFields
  mortgage *mortgageFields
  bonds *bondsFields
  adFv *adFvFields
  adPv *adPvFields
  adCp *adCpFields
  adEpp *adEppFields
  oaCp *oaCpFields
  oaEpp *oaEppFields
  oaFv *oaFvFields
  oaGa *oaGaFields
  oaInterestRate *oaInterestRateFields
  oaPerpetuity *oaPerpetuityFields
  oaPv *oaPvFields
  siAccurate *siAccurateFields
  siBankers *siBankersFields
  siOrdinary *siOrdinaryFields
}

//Called from main.go before requests are accepted.
func SetupDirStructure(dir string) {
  mainDir = filepath.Join(dir, mainDir)
  dirErr, err := osu.CreateDirs(0o077, 0o777, mainDir)
  if err != nil {
    panic("Cannot create directory '" + dirErr + "': " + err.Error())
  }
}

/***
Even though each username has its own exclusive set of files, multiple users can share the same username. In this scenario,
there is a possibility that two or more users may try to modify a file at the same time thereby corrupting the file. In order
to prevent file corruption, the file is protected with a lock that allows a single writer (exclusive write) or multiple readers
(share reads).
***/
func readFields(filePath string) ([]byte, error) {
  exists, err := osu.CheckFileExists(filePath)
  if exists {
    obj, err := osu.ReadAllShareLock1(filePath, os.O_RDONLY, 0o400)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Couldn't open file %s: ", filePath) + err.Error())
    }
    return obj, nil
  } else {
    return nil, err
  }
}

type UserSession struct{
  Fields *fields
  RefCount int
  /***
  If a user closes its browser without logging out, or goes inactive, the data structures will stay stuck in memory forever
  because DeleteSessionDataPerUser is never triggered by an HTTP action. This will lead to a slow, permanent memory leak.

  The LastAccessed field is used to avoid this issue; a background goroutine wakes up every few minutes, locks the manager,
  checks for expired users, and evicts them.
  ***/
  LastAccessed time.Time  //Track when this session was last used.
}

/***
loadFieldsFromDisk reads a JSON file from disk, unmarshals it into a struct populated with defaults, and returns the result.
If reading or unmarshaling fails, it safely falls back to defaults.
***/
func loadFieldsFromDisk[T any](dir1, dir2, filename, correlationId string, defaultFields *T) *T {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)  //Creation of dirs is protected by lock.
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  fullFilePath := filepath.Join(dir, filename)
  obj, err := readFields(fullFilePath)  //Protected with multiple reads/single writer lock.
  if err != nil {
    logger.LogError(fmt.Sprintf("Failed to read fields from file, continue with defaults: %+v", err), correlationId)
    return defaultFields
  }
  //
  if obj == nil {
    logger.LogInfo(fmt.Sprintf("File %s does not exist; continue with defaults.", fullFilePath), correlationId)
    return defaultFields
  }
  /***
  When a file is empty, the readFields function successfully returns a valid slice, but it contains zero bytes. Checking the
  length ensures parsing only files that actually contain data.
  ***/
  if len(obj) == 0 {  //Check if the file contains no data (empty).
    logger.LogInfo(fmt.Sprintf("File %s is empty, continue with defaults.", fullFilePath), correlationId)
    return defaultFields
  }
  /***
  Create a fresh scratchpad copy, pre-loaded with the default values. If the JSON on disk contains only a few fields, the
  missing fields will be set to defaults.
  ***/
  resultCopy := *defaultFields
  //json.Unmarshal can partially modify the destination variable even if it returns an error.
  err = json.Unmarshal(obj, &resultCopy)
  if err != nil {
    // Write error, but continue with default values.
    logger.LogWarning(fmt.Sprintf("Corrupt JSON structure inside %s; continue with defaults. Error: %+v", fullFilePath, err), correlationId)
    return defaultFields
  }
  return &resultCopy
}

/***
getFieldsGeneric handles the read-locking, map validation, and shallow-copying for any user field type stored within the global
currentFields map.
***/
func getFieldsGeneric[T any](userName string, selector func(*UserSession) *T) *T {
  //Acquire a read lock before opening the map structure.
  currentFieldsLock.RLock()
  //Ensure the lock is released when this function finishes executing.
  defer currentFieldsLock.RUnlock()
  //Safely check if the user exists in the map to prevent nil-pointer panics.
  userSession, exists := currentFields[userName]
  if !exists || userSession == nil || userSession.Fields == nil {
    return nil  //Safely return nil if the session or internal fields don't exist.
  }
  //Use the selector callback to extract the specific pointer field (e.g., .adCp or .adEpp).
  fieldPtr := selector(userSession)
  if fieldPtr == nil {
    return nil
  }
  /***
  One Important Caveat (Shallow vs. Deep Copy)
  A shallow copy is 100% safe if fields only contains primitive types (like string, int, bool, float64). However, if fields
  contains internal reference types like pointers, slices, or maps, a shallow copy will copy the slice/map header, meaning
  the new pointer will still point to the same underlying data array.
  1. *userSession.Fields.adCp dereferences the pointer to read the literal string values.
  2. The outer '&' takes that copied value and spins up a brand new memory block on the heap.
  3. Since strings are immutable, this new pointer is 100% independent from the global map.
  ***/
  isolatedCopy := *fieldPtr
  return &isolatedCopy
}

/***
Because this function modifies the global currentFields map and triggers file/directory creation routines, running it
without protection will guarantee race conditions and application crashes in production.
Called from VerifyLogin (WfFinPages.go) after a user has been verified.
***/
func AddSessionDataPerUser(userName, correlationId string) {
  /***
  Because the LockUser guard spans the entire duration of the fields initialization block, it is physically impossible for a
  duplicate login request to slide past or exit early with uninitialized data; i.e., no two goroutines can execute
  AddSessionDataPerUser for the same username at the same time.
  ***/
  unlock := UserLockMgr.LockUser(userName)
  defer unlock()
  //Check if a session already exists.
  currentFieldsLock.Lock()
  session, exists := currentFields[userName]
  if exists {
    session.RefCount++
    session.LastAccessed = time.Now()
    currentFieldsLock.Unlock()
    return
  }
  //Release the lock while doing heavy file/directory creation operations so other users aren't blocked from accessing their sessions.
  currentFieldsLock.Unlock()
  fd := &fields{  //The newXxxFields functions are thread-safe.
    miscellaneous: newMiscellaneousFields(mainDir, userName, correlationId),
    mortgage: newMortgageFields(mainDir, userName, correlationId),
    bonds: newBondsFields(mainDir, userName, correlationId),
    adFv: newAdFvFields(mainDir, userName, correlationId),
    adPv: newAdPvFields(mainDir, userName, correlationId),
    adCp: newAdCpFields(mainDir, userName, correlationId),
    adEpp: newAdEppFields(mainDir, userName, correlationId),
    oaCp: newOaCpFields(mainDir, userName, correlationId),
    oaEpp: newOaEppFields(mainDir, userName, correlationId),
    oaFv: newOaFvFields(mainDir, userName, correlationId),
    oaGa: newOaGaFields(mainDir, userName, correlationId),
    oaInterestRate: newOaInterestRateFields(mainDir, userName, correlationId),
    oaPerpetuity: newOaPerpetuityFields(mainDir, userName, correlationId),
    oaPv: newOaPvFields(mainDir, userName, correlationId),
    siAccurate: newSiAccurateFields(mainDir, userName, correlationId),
    siBankers: newSiBankersFields(mainDir, userName, correlationId),
    siOrdinary: newSiOrdinaryFields(mainDir, userName, correlationId),
  }
  //Lock the main map briefly to save the new pointer.
  currentFieldsLock.Lock()
  currentFields[userName] = &UserSession{
    Fields: fd,
    RefCount: 1,
    LastAccessed: time.Now(),  //Track the exact creation timestamp.
  }
  currentFieldsLock.Unlock()
}

func DeleteSessionDataPerUser(userName, correlationId string) {

logger.LogInfo("*********** DeleteSessionDataPerUser ***********", correlationId)


  unlock := UserLockMgr.LockUser(userName)
  defer unlock()
  /***
  Using defer on a global lock inside a fast map operation forces the global currentFieldsLock to stay locked slightly longer
  than necessary while Go processes the function exit. Since deleting an element from a map takes only a few nanoseconds,
  manually unlock the map lock to free up other users immediately.
  Lock the global session map.
  ***/
  currentFieldsLock.Lock()
  /***
  Under the current architecture, this check is logically impossible to hit during standard application flows. But if a rogue
  event hits an uninitialized session, exit safely instead of panicking.
  ***/
  session, exists := currentFields[userName]
  if !exists {
    currentFieldsLock.Unlock()
    return
  }
  //Decrement the session reference counter.
  session.RefCount--
  //Clean up the session data from memory only when active browser tabs/requests drop to zero.
  if session.RefCount < 1 {
    delete(currentFields, userName)
    logger.LogInfo(fmt.Sprintf("(DeleteSessionDataPerUser) Deleted user %s.", userName), correlationId)
  }
  currentFieldsLock.Unlock()
}

/***
If a user closes their browser without logging out, or goes inactive, the data structures will stay stuck in memory forever because
DeleteSessionDataPerUser is never triggered by an HTTP action. This will lead to a slow, permanent memory leak.

To prevent this memory leak on long-running servers, the janitor must check expiration before it invokes the lock manager, and only
escalate to the user lock if eviction is actually required.
***/
func StartSessionJanitor(timeoutDuration time.Duration, checkInterval time.Duration, correlationId string) {
  //Initialize the permanent, memory-safe cleanup thread.
  go func() {
    logger.LogInfo("Starting the Janitor goroutine StartSessionJanitor", correlationId)
    //Create a time.Ticker to execute code repeatedly at a fixed, regular interval.
    ticker := time.NewTicker(checkInterval)
    //Always ensure to stop the ticker to free up resources.
    defer ticker.Stop()
    //It automatically blocks and loops indefinitely, executing the body of the loop every time the ticker sends a new
    //value down its channel.
    for range ticker.C {



logger.LogInfo("*********** StartSessionJanitor ***********", correlationId)




      now := time.Now()  //Return the current local time as a time.Time object.
      var expiredUsers []string
      //Read-Lock the global map to identify expired users.
      currentFieldsLock.RLock()
      for userName, session := range currentFields {
        if now.Sub(session.LastAccessed) > timeoutDuration {
          expiredUsers = append(expiredUsers, userName)
        }
      }
      currentFieldsLock.RUnlock()
      //Evict only the expired users.
      for _, userName := range expiredUsers {
        /***
        Check status using a cheap READ-lock first. This keeps the user path unblocked while the janitor decides if it's
        worth waiting for the user lock.
        ***/
        currentFieldsLock.RLock()
        session, stillExists := currentFields[userName]
        if !stillExists || !(timeoutDuration > time.Since(session.LastAccessed)) {
          currentFieldsLock.RUnlock()
          continue  //Skip entirely! No LockUser call, no refCount adjustments.
        }
        currentFieldsLock.RUnlock()
        //Now it is highly likely we actually need to evict. Acquire the user lock.
        unlock := UserLockMgr.LockUser(userName)
        //Secure the full write-lock to safely mutate the global session map.
        currentFieldsLock.Lock()
        //Check to see if it was removed while waiting for the lock.
        if currentSession, exists := currentFields[userName]; exists {
          /***
          There is a blocking call (UserLockMgr.LockUser). The thread might sleep for several milliseconds waiting for
          a busy user lock, making the original time.Now snapshot heavily outdated by the time the lock is acquired. Hence,
          check if the session is currently expired relative to the immediate present when the lock is finally acquired.
          Double check to verify that an action was not triggered while waiting for the lock.
          ***/
          if time.Since(currentSession.LastAccessed) > timeoutDuration {
            delete(currentFields, userName)
            logger.LogInfo(fmt.Sprintf("(StartSessionJanitor) Deleted user %s.", userName), correlationId)
          }
        }
        currentFieldsLock.Unlock()
        unlock()
      }
    }
  }()
}
