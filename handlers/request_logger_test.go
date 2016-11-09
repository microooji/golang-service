package handlers

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUri(t *testing.T) {
	cases := map[string]struct {
		req      *http.Request
		expected string
	}{
		"basic": {
			newRequest("GET", "http://example.com"),
			"/",
		},
		"http2": {
			&http.Request{
				Method:     "CONNECT",
				Proto:      "HTTP/2.0",
				ProtoMajor: 2,
				ProtoMinor: 0,
				URL:        &url.URL{Host: "www.example.com:443"},
				Host:       "www.example.com:443",
				RemoteAddr: "192.168.100.5",
			},
			"www.example.com:443",
		},
		"with params": {
			newRequest("GET", "http://example.com/path?query=value"),
			"/path?query=value",
		},
		"without params": {
			newRequest("GET", "http://example.com/path"),
			"/path",
		},
	}

	for k, tc := range cases {
		url := *tc.req.URL
		assert.Equal(t, tc.expected, parseURI(tc.req, url), "test %s", k)
	}
}

func TestUriPath(t *testing.T) {
	cases := map[string]struct {
		req      *http.Request
		expected string
	}{
		"basic": {
			newRequest("GET", "http://example.cpom"),
			"/",
		},
		"http2": {
			&http.Request{
				Method:     "CONNECT",
				Proto:      "HTTP/2.0",
				ProtoMajor: 2,
				ProtoMinor: 0,
				URL:        &url.URL{Host: "www.example.com:443"},
				Host:       "www.example.com:443",
				RemoteAddr: "192.168.100.5",
			},
			"www.example.com:443",
		},
		"with params": {
			newRequest("GET", "http://example.com/path?query=value"),
			"/path",
		},
		"without params": {
			newRequest("GET", "http://example.com/path"),
			"/path",
		},
	}

	for k, tc := range cases {
		url := *tc.req.URL
		assert.Equal(t, tc.expected, uriPath(tc.req, url), "test %s", k)
	}
}

func TestGetUserIP(t *testing.T) {
	cases := map[string]struct {
		req      *http.Request
		expected interface{}
	}{
		"undefined": {
			newRequest("GET", "http://example.com"),
			(net.IP)(nil),
		},
	}

	for k, tc := range cases {
		ip, _ := getUserIP(tc.req)
		assert.Equal(t, tc.expected, ip, "test %s", k)
	}
}
