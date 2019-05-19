package goquest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	uri, err := url.Parse(rawurl)
	if err != nil {
		log.Println("Goquest:", err)
	}
	return &Goquest{
		url: rawurl,
		request: &http.Request{
			URL:    uri,
			Method: method,
			Header: http.Header{
				"User-Agent": []string{"Goquest"},
			},
		},
		timeout: 30,
	}
}

func (g *Goquest) Query() (*Goquest, error) {
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

// Request
// SetUserAgent sets User-Agent header field
func (g *Goquest) SetUserAgent(ua string) *Goquest {
	g.request.Header.Set("User-Agent", ua)
	return g
}

// Response
func (g *Goquest) Byte() []byte {
	defer g.response.Body.Close()
	if g.response == nil || g.response.Body == nil {
		fmt.Println("Goquest: Response Or Body Is Nil")
		return []byte{}
	}
	if g.body == nil {
		g.body, _ = ioutil.ReadAll(g.response.Body)
	}
	return g.body
}

func (g *Goquest) String() string {
	return string(g.Byte())
}

func (g *Goquest) Json(v interface{}) error {
	data := g.Byte()
	return json.Unmarshal(data, v)
}

func (g *Goquest) StatusCode() int {
	if g.response == nil {
		return 0
	}
	return g.response.StatusCode
}
