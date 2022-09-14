package drivers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

type SlideCurtain interface {
	SetPosition(position float64) bool
}

type slideCurtain struct {
	deviceId   string
	ip         string
	nonceCount int
}

func CreateSlideCurtain(ip string, deviceId string) SlideCurtain {
	return &slideCurtain{
		deviceId:   deviceId,
		ip:         ip,
		nonceCount: 0,
	}
}

type slidePositionRequest struct {
	Position float64 `json:"pos"`
}

func (sc *slideCurtain) SetPosition(position float64) bool {
	if position > 1 {
		position = 1
	} else if position < 0 {
		position = 0
	}

	body, _ := json.Marshal(slidePositionRequest{Position: position})

	err := sc.postWithDigestAuth("/rpc/Slide.SetPos", body)

	if err != nil {
		return false
	}

	return true
}

func (sc *slideCurtain) postWithDigestAuth(uri string, body []byte) error {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	response, err := sc.sendRequest("POST", uri, headers, []byte{})

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("unexpected status code, expected %d (with auth challenge), got %d", http.StatusUnauthorized, response.StatusCode)
	}

	authChallenge := response.Header.Get("www-authenticate")

	if authChallenge == "" {
		return errors.New("no auth challenge present in initial response from device (expected www-authenticate header)")
	}

	qopRegex, _ := regexp.Compile("qop=\"(.*?)\"")
	realmRegex, _ := regexp.Compile("realm=\"(.*?)\"")
	nonceRegex, _ := regexp.Compile("nonce=\"(.*?)\"")

	qopMatch := qopRegex.FindAllStringSubmatch(authChallenge, 1)
	realmMatch := realmRegex.FindAllStringSubmatch(authChallenge, 1)
	nonceMatch := nonceRegex.FindAllStringSubmatch(authChallenge, 1)

	if qopMatch == nil || realmMatch == nil || nonceMatch == nil {
		return fmt.Errorf("auth challenge is missing realm or nonce (was: %s)", authChallenge)
	}

	qop := qopMatch[0][1]
	realm := realmMatch[0][1]
	nonce := nonceMatch[0][1]

	ha1 := hashMD5(fmt.Sprintf("%s:%s:%s", "user", realm, sc.deviceId))
	ha2 := hashMD5(fmt.Sprintf("%s:%s", "POST", uri))

	sc.nonceCount++

	nc := fmt.Sprintf("%08d", sc.nonceCount)
	cnonce := fmt.Sprintf("%08s", randomNonce())

	res := hashMD5(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, nonce, nc, cnonce, qop, ha2))

	headers["Authorization"] = fmt.Sprintf("Digest username=\"user\", realm=\"%s\", nonce=\"%s\", uri=\"%s\", response=\"%s\", qop=\"%s\", nc=\"%s\", cnonce=\"%s\"", realm, nonce, uri, res, qop, nc, cnonce)

	apiResponse, err := sc.sendRequest("POST", uri, headers, body)

	if err != nil {
		fmt.Println(err)
		fmt.Println(apiResponse)
	}

	return nil
}

func (sc *slideCurtain) sendRequest(method string, uri string, headers map[string]string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("http://%s%s", sc.ip, uri)

	request, _ := http.NewRequest(method, url, bytes.NewBuffer(body))

	for name, value := range headers {
		request.Header.Set(name, value)
	}

	return http.DefaultClient.Do(request)
}

func hashMD5(src string) string {
	ha1 := md5.Sum([]byte(src))

	return fmt.Sprintf("%x", ha1)
}

func randomNonce() string {
	n := make([]byte, 4)

	rand.Seed(time.Now().UnixNano())
	rand.Read(n)

	return hex.EncodeToString(n)
}
