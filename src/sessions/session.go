package sessions

import (
  //"sync"
  "time"
)

type Session struct {
  Username string
  Expiry time.Time
//  lock sync.Mutex  //Protect session
}

//Create a struct that models the structure of a user, both in the request body, and in the DB
// type Credentials struct {
//   Password string `json:"password"`
//   Username string `json:"username"`
// }
type Credentials struct {
  Password string
  Username string
}

//Store information about each user's session.
var Sessions = map[string]Session{}  //key: sessionToken, value: session

var Users = map[string]string{}  //key: username, value: password

//Determine if a session has expired.
func (s *Session) IsExpired() bool {
  return s.Expiry.Before(time.Now())
}



