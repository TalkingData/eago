package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var httpCli = &http.Client{}

const (
	BASE_ADDRESS = "http://localhost:8080"
	CONTENT_TYPE = "application/json"
	TOKEN        = "my_token"
)

func TestApiNewNode(t *testing.T) {
	t.Log("TestApiNewNode")
	uri := BASE_ADDRESS + "/flow/nodes"
	req, err := http.NewRequest(
		"POST",
		uri,
		strings.NewReader(`{"name":"test_node_1","category":1,"visible_fields":"","editable_fields":""}`),
	)
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Content-Type", CONTENT_TYPE)
	req.Header.Set("token", TOKEN)
	resp, err := httpCli.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}

	t.Log(string(body))
}
