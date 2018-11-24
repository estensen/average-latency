package main

import "testing"



func TestLatency(t *testing.T) {
	url := "google.no"
	firstLatency, lastLatency := GetLatency(url)

	if firstLatency > lastLatency {
		t.Error("Latency for receiving all bytes is less than for the first byte")
	}
}
