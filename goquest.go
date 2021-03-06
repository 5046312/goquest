package goquest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Goquest struct {
	url      string
	params   map[string][]string
	request  *http.Request
	response *http.Response
	body     []byte
	timeout  time.Duration
}

func newGoquest(rawurl, method string) *Goquest {
	//var resp http.Response
	rUrl, err := url.Parse(rawurl)
	if err != nil {
		log.Println("Goquest:", err)
	}
	return &Goquest{
		url:    rawurl,
		params: make(map[string][]string),
		request: &http.Request{
			URL:    rUrl,
			Method: method,
			Header: http.Header{
				"User-Agent": []string{"Goquest"},
			},
		},
		timeout: 30,
	}
}

func (g *Goquest) Query() (*Goquest, error) {
	// encodeGetParams
	g.encodeGetParams()
	// Post Params
	g.encodePostParams()

	client := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * g.timeout,
		},
	}
	resp, err := client.Do(g.request)
	if err != nil {
		return nil, err
	}
	g.response = resp
	return g, nil
}

func Get(url string) *Goquest {
	return newGoquest(url, http.MethodGet)
}

func Post(url string) *Goquest {
	return newGoquest(url, http.MethodPost)
}

func (g *Goquest) Param(key, value string) *Goquest {
	if param, ok := g.params[key]; ok {
		g.params[key] = append(param, value)
	} else {
		g.params[key] = []string{value}
	}
	return g
}

// Get Method: Params To Url
func (g *Goquest) encodeGetParams() {
	// 仅限Get请求,将params转成get参数
	if g.request.Method == http.MethodGet && len(g.params) != 0 {
		params := url.Values(g.params)
		paramString := params.Encode()
		oUrl := g.url
		if !strings.Contains(oUrl, "?") {
			// Not Contain `?`
			oUrl += "?" + paramString
		} else {
			oUrl += "&" + paramString
		}
		uri, err := url.Parse(oUrl)
		if err == nil {
			g.request.URL = uri
		}
	}
}

// Post or other Method: Set Body
func (g *Goquest) encodePostParams() {
	if g.request.Method != http.MethodGet && g.request.Body == nil && len(g.params) > 0 {
		params := url.Values(g.params)
		paramString := params.Encode()
		bf := bytes.NewBufferString(paramString)
		g.request.Body = ioutil.NopCloser(bf)
		g.request.ContentLength = int64(len(paramString))
		g.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
}

// JsonBody
func (g *Goquest) JsonBody(data interface{}) (*Goquest, error) {
	if g.request.Body == nil && data != nil {
		dataByte, err := json.Marshal(data)
		if err != nil {
			return g, err
		}
		g.request.Body = ioutil.NopCloser(bytes.NewReader(dataByte))
		g.request.ContentLength = int64(len(dataByte))
		g.request.Header.Set("Content-Type", "application/json")
	}
	return g, nil
}

// Request
// Set Header Map
func (g *Goquest) SetHeader(key, value string) *Goquest {
	g.request.Header.Set(key, value)
	return g
}

// SetUserAgent sets User-Agent header field
func (g *Goquest) SetUserAgent(ua string) *Goquest {
	g.SetHeader("User-Agent", ua)
	return g
}

// Response
func (g *Goquest) Byte() []byte {
	if g.StatusCode() == 0 {
		fmt.Println("Goquest: You May Have Forgotten To Call The `Query`")
		return []byte{}
	}
	g.body, _ = ioutil.ReadAll(g.response.Body)
	g.response.Body.Close()
	return g.body
}

func (g *Goquest) String() string {
	return string(g.Byte())
}

func (g *Goquest) Json(v interface{}) error {
	return json.Unmarshal(g.Byte(), v)
}

func (g *Goquest) StatusCode() int {
	if g.response == nil || g.response.Body == nil {
		return 0
	}
	return g.response.StatusCode
}
