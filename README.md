Tokbox Golang
=============

A wrapper for the Tokbox http api.

Install
-------

```shell
go get github.com/cioc/tokbox.git
```

Usage
-----

```go
import "github.com/cioc/tokbox"

//setup the api to use your credentials
tb := tokbox.New("<my api key>","<my secret key>")

//create a session
session, err := tb.NewSession("", true) //no location, peer enabled

//create a token
token, err := session.token("publisher", "", 86400) //type publisher, no connection data, expire in 24 hours

```


See the tests for more detailed examples.


Sources: 
--------
(These were both very helpful)

http://www.tokbox.com/forums/other-server-api/session-and-token-generation-in-go-gae-golang-org-t17644
http://www.tokbox.com/blog/generating-tokens-without-server-side-sdk/

