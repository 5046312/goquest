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
	}
}
