package goquest

import "testing"

func TestGet(t *testing.T) {
	quest, err := Get("http://httpbin.org/get").Query()
	if err == nil {
		t.Log(quest.String())
		t.Log(quest.Byte())

		type R struct {
			Origin string `json:"origin"`
		}
		json := &R{}
		err := quest.Json(json)
		if err != nil {
			t.Error(err)
		}
		t.Log(json.Origin)
	} else {
		t.Error(err)
	}
}

func TestGetErr(t *testing.T) {
	quest, err := Get("http://err.xxx.com/ppp").Query()
	if err != nil {
		t.Log(err)
	}
	//t.Log(quest.StatusCode())
	t.Log(quest)
}

func TestGoquest_SetUserAgent(t *testing.T) {
	quest, err := Get("http://httpbin.org/get").SetUserAgent("pppp").Query()
	if err == nil {
		t.Log(quest.String())
	}
}
