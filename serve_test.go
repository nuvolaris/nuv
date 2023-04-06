package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	handler := nuvServerHandler(joinpath("tests/olaris/", WebDir))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	newreq := func(method, url string, body io.Reader) *http.Request {
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		name string
		r    *http.Request
	}{
		{name: "1: testing get", r: newreq("GET", ts.URL+"/", nil)},
		// {name: "2: testing post", r: newreq("POST", ts.URL+"/", nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.r)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
			}
			resp.Body.Close()
		})
	}

}
