package tokbox

import (
  "testing"
  "log"
)

const key = "<your key here>"
const secret = "<your secret here>"

func TestNewSession(t *testing.T) {
  tokbox := New(key, secret)
  session, err := tokbox.NewSession("", true)
  if err != nil {
    log.Fatal(err)
    t.FailNow()
  }
  log.Println(session)
}

func TestToken(t *testing.T) {
  tokbox := New(key, secret)
  session, err := tokbox.NewSession("", true)
  if err != nil {
    log.Fatal(err)
    t.FailNow()
  }
  token, err := session.Token("publisher", "", 86400)
  if err != nil {
    log.Fatal(err)
    t.FailNow()
  }
  log.Println(token)
}
