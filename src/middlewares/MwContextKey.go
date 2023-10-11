package middlewares

import (
  "context"
)

/***
To avoid context keys collisions, a best practice is to create an unexported custom type.
***/
type ctxKey string
/***
The correlationIdKey constant is unexported. Hence, there's no risk that another package using the
same context could override the value that is already set. Even if another package creates the same
correlationIdKey based on a ctxKey type as well, it will be a different key.
***/
const correlationIdKey ctxKey = "correlationIdKey"
const sessionTokenKey ctxKey = "sessionTokenKey"

type MwContextKey struct{}

/***
Packages that define a Context key should provide type-safe accessors for the values stored using
that key.
***/
func (ck MwContextKey) GetCorrelationId(ctx context.Context) (cid string, ok bool) {
  //A type assertion is an operation applied to an interface value.
  cid, ok = ctx.Value(correlationIdKey).(string)
  return
}

func (ck MwContextKey) GetSessionToken(ctx context.Context) (sessionToken string, ok bool) {
  //A type assertion is an operation applied to an interface value.
  sessionToken, ok = ctx.Value(sessionTokenKey).(string)
  return
}
