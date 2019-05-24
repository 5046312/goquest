package goquest

import "testing"

func TestGetString(t *testing.T) {
	quest, err := Get("http://httpbin.org/get").Query()
	if err == nil {
		t.Log(quest.String())
	} else {
		t.Error(err)
	}
}

func TestGetJson(t *testing.T) {
	quest, err := Get("http://httpbin.org/get").Query()
	if err == nil {
		type R struct {
			Origin string `json:"origin"`
		}
		json := &R{}
		err := quest.Json(json)
		if err != nil {
			t.Error(err)
		}
		t.Log(json.Origin)
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
func TestGoquest_SetGetParams(t *testing.T) {
	quest, err := Get("http://httpbin.org/get").Param("a", "1").Param("b", "2").Param("c", "3").Query()
	if err == nil {
		t.Log(quest.String())
	}

	quest, err = Get("http://httpbin.org/get?get=12345").Param("a", "1").Param("b", "2").Param("c", "3").Query()
	if err == nil {
		t.Log(quest.String())
	}
}
