package helper

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/microcosm-cc/bluemonday"
)

var policy = bluemonday.NewPolicy().AllowElements("b", "br", "p")

var Client *http.Client

func init() {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RetryMax = 2
	Client = retryClient.StandardClient()
	Client.Timeout = time.Duration(60) * time.Second
}

// Println  ...
func Println(i interface{}) {
	fmt.Println("-- begin --")
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
	fmt.Println("-- end --")
}

// Fetch fetch a url
func Fetch(ctx context.Context, req *http.Request) (*http.Response, error) {
	// set default referer
	if req.Header.Get("Referer") == "" {
		req.Header.Set("Referer", req.URL.String())
	}

	if req.Header.Get("Referer") == "noref" {
		req.Header.Del("Referer")
	}

	// default user-agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	}

	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	// check status 200
	if res.StatusCode != 200 {
		res.Body.Close()
		return nil, errors.New(req.URL.String() + ":" + res.Status)
	}
	return res, nil
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var digits = []rune("0123456789")

func RandDigitString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func UnsafeHTML(s string) string {
	return policy.Sanitize(strings.ReplaceAll(html.UnescapeString(s), "\n", "<br>"))
}

// Reverse reverse a string
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func Decode(str string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(Reverse(str))
	return string(data), err
}

func Encode(str string) string {
	return Reverse(base64.RawURLEncoding.EncodeToString([]byte(str)))
}

// true of s1 == s2
func EqualStringSlice(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// Contains return true if `s“ in `list`
func Contains(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
