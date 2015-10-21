package main

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestGetStubStatus(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 7 8 19")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts1.Close()

	u1, _ := url.Parse(ts1.URL)
	nb1 := Nginxbeat{url: u1}
	s1, _ := nb1.getStubStatus()

	assert.Equal(t, s1["active"], 1)
	assert.Equal(t, s1["accepts"], 7)
	assert.Equal(t, s1["handled"], 8)
	assert.Equal(t, s1["requests"], 19)
	assert.Equal(t, s1["reading"], 0)
	assert.Equal(t, s1["writing"], 1)
	assert.Equal(t, s1["waiting"], 2)

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 7 8 19")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts2.Close()

	u2, _ := url.Parse(ts2.URL)
	nb2 := Nginxbeat{url: u2}
	s2, _ := nb2.getStubStatus()

	assert.Equal(t, s2["active"], -1)
	assert.Equal(t, s2["accepts"], 7)
	assert.Equal(t, s2["handled"], 8)
	assert.Equal(t, s2["requests"], 19)
	assert.Equal(t, s2["reading"], 0)
	assert.Equal(t, s2["writing"], 1)
	assert.Equal(t, s2["waiting"], 2)

	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts3.Close()

	u3, _ := url.Parse(ts3.URL)
	nb3 := Nginxbeat{url: u3}
	s3, _ := nb3.getStubStatus()

	assert.Equal(t, s3["active"], 1)
	assert.Equal(t, s3["accepts"], -1)
	assert.Equal(t, s3["handled"], -1)
	assert.Equal(t, s3["requests"], -1)
	assert.Equal(t, s3["reading"], 0)
	assert.Equal(t, s3["writing"], 1)
	assert.Equal(t, s3["waiting"], 2)

	ts4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 7 8 19")
		fmt.Fprintln(w, "")
	}))
	defer ts4.Close()

	u4, _ := url.Parse(ts4.URL)
	nb4 := Nginxbeat{url: u4}
	s4, _ := nb4.getStubStatus()

	assert.Equal(t, s4["active"], 1)
	assert.Equal(t, s4["accepts"], 7)
	assert.Equal(t, s4["handled"], 8)
	assert.Equal(t, s4["requests"], 19)
	assert.Equal(t, s4["reading"], -1)
	assert.Equal(t, s4["writing"], -1)
	assert.Equal(t, s4["waiting"], -1)
}