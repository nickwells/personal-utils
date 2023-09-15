package main

import (
	"testing"
	"time"
)

func TestClientConnErrors(t *testing.T) {
	var testHosts = []struct {
		hostAddr string
		desc     string
		expectOK bool
	}{
		{"192.168.1.97:22", "raspberry pi", true},
		{"192.168.1.64:22", "printer - no ssh", false},
		{"192.168.1.65:22", "nonesuch", false},
	}

	for _, hd := range testHosts {
		_, err := makeClientConn(hd.hostAddr)
		t.Log(time.Now(), hd.desc, err, "\n")
		if err != nil && hd.expectOK {
			t.Error("couldn't make the connection to: " + hd.hostAddr +
				" - " + hd.desc + ". err: " + err.Error())
		} else if err == nil && !hd.expectOK {
			t.Error("no error detected while connecting to: " + hd.hostAddr +
				" - " + hd.desc)
		}

	}
}
