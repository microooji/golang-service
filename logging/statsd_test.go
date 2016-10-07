// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package logging

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/DataDog/datadog-go/statsd"
    "../nettest"
    "time"
    "os"
    "net/http"
)

func TestStatsdLogging(t *testing.T) {
    done := make(chan string)
    addr, sock, srvWg := nettest.CreateServer(t, "udp", "localhost:", done)
    defer srvWg.Wait()
    defer os.Remove(addr.String())
    defer sock.Close()

	client, err := statsd.New(addr.String())
	if err != nil {
		t.Fatal(err)
	}
    client.Tags = append(client.Tags, "test")
    client.Namespace = "service.logging.live."

    loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}
	ts := time.Date(1983, 05, 26, 3, 30, 45, int((736 * time.Millisecond).Nanoseconds()), loc)

	// A typical request with an OK response
	req := newRequest("GET", "http://example.com")

    dur, err := time.ParseDuration("0.302s")
    if (err != nil) {
        t.Fatal(err)
    }

    writeStatsdLog(client, req, *req.URL, ts, dur, http.StatusOK, 100)
    expected := []string{
        "service.logging.live.request.response_time:302.000000|ms|#test,endpoint:/,statusCode:200,method:GET",
        "service.logging.live.request.count:1|c|#test,endpoint:/,statusCode:200,method:GET",
    }

    for _, message := range expected {
        assert.Equal(t, message, <-done)
    }

    ts = time.Date(1983, 05, 26, 3, 30, 45, int((123 * time.Millisecond).Nanoseconds()), loc)
    req = newRequest("POST", "http://example.com/path/here")

    dur, err = time.ParseDuration("0.102s")
    if (err != nil) {
        t.Fatal(err)
    }

    expected = []string{
        "service.logging.live.request.response_time:102.000000|ms|#test,endpoint:/path/here,statusCode:200,method:POST",
        "service.logging.live.request.count:1|c|#test,endpoint:/path/here,statusCode:200,method:POST",
    }

    writeStatsdLog(client, req, *req.URL, ts, dur, http.StatusOK, 100)

    for _, message := range expected {
        assert.Equal(t, message, <-done)
    }
}
