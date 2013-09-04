package tokbox

import (
  "bytes"

  "net/url"
  "net/http"

  "io/ioutil"

  "encoding/xml"
  "encoding/base64"

  "crypto/hmac"
  "crypto/sha1"

  "fmt"
  "strings"
  "time"
  "math/rand"
)

const apiHost = "https://api.opentok.com/hl"
const apiSession = "/session/create"

type Tokbox struct {
  apiKey string
  partnerSecret string
}

type Session struct {
  SessionId string      `xml:"session_id"`
  PartnerId string      `xml:"partner_id"`
  CreateDt string       `xml:"create_dt"`
  SessionStatus string  `xml:"session_status"`
  t *Tokbox
}

//private - only for parsing xml purposes
type sessions struct {
  Sessions []Session `xml:"Session"`
}

func New(apikey, partnerSecret string) (*Tokbox) {
  return &Tokbox{apikey, partnerSecret}
}

//remember expiration = 86400 is 24 hours
func (s *Session) Token(role string, connectionData string, expiration int64) (string, error) {
  now := time.Now().Unix()

  dataVals := url.Values{}
  dataVals.Add("session_id", s.SessionId)
  dataVals.Add("create_time", fmt.Sprintf("%d", now))
  dataVals.Add("expire_time", fmt.Sprintf("%d", now + expiration))
  dataVals.Add("role", role)
  dataVals.Add("connection_data", connectionData)
  dataVals.Add("nonce", fmt.Sprintf("%d", rand.Intn(999999)))
  dataStr := dataVals.Encode()

  h := hmac.New(sha1.New, []byte(s.t.partnerSecret))
  n, err := h.Write([]byte(dataStr))
  if err != nil {
    return "", err
  }
  if n != len(dataStr) {
    return "", fmt.Errorf("hmac not enough bytes written %d != %d", n, len(dataStr))
  }

  tk := url.Values{}
  tk.Add("partner_id", s.t.apiKey)
  tk.Add("sig", fmt.Sprintf("%x:%s", h.Sum(nil), dataStr))

  var buf bytes.Buffer
  encoder := base64.NewEncoder(base64.StdEncoding, &buf)
  encoder.Write([]byte(tk.Encode()))
  encoder.Close()
  return fmt.Sprintf("T1==%s", buf.String()), nil
}

func (t *Tokbox) NewSession(location string, p2p bool) (*Session, error) {
  params := url.Values{}
  if len(location) > 0 {
    params.Add("location", location)
  }
  if p2p {
    params.Add("p2p.preference", "enabled")
  } else {
    params.Add("p2p.preference", "disabled")
  }
  req, err := http.NewRequest("POST", apiHost + apiSession, strings.NewReader(params.Encode()))
  if err != nil {
    return &Session{}, err
  }
  authHeader := t.apiKey + ":" + t.partnerSecret
  req.Header.Add("X-TB-PARTNER-AUTH", authHeader)
  client := http.Client{}
  res, err := client.Do(req)
  defer res.Body.Close()
  if err != nil {
    return &Session{}, err
  }
  if res.StatusCode != 200 {
    return &Session{}, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
  }
  b, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return &Session{}, err
  }
  var s sessions
  err = xml.Unmarshal(b, &s)
  if err != nil {
    return &Session{}, err
  }
  if len(s.Sessions) < 1 {
    return &Session{}, fmt.Errorf("tokbox did not return a session")
  }
  o := s.Sessions[0]
  o.t = t
  return &o, nil
}
