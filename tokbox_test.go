package tokbox

import (
	"log"
	"testing"
)

const key = "<your api key here>"
const secret = "<your partner secret here>"

func TestToken(t *testing.T) {
	tokbox := New(key, secret)
	session, err := tokbox.NewSession("", true)
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}
	log.Println(session)
	token, err := session.Token("", "", -1) //defaults to publisher, no connection data and expires in 24 hours
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}
	log.Println(token)
}
