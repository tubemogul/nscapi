package main

import (
	"github.com/PurpureGecko/go-lfc"
	nsca "github.com/tubemogul/nscatools"
	"testing"
)

func TestQueueData(t *testing.T) {
	q = lfc.NewQueue()
	if q.Len() > 0 {
		t.Error("Queue not empty on start")
	}

	p := &nsca.DataPacket{HostName: "host01", Service: "service foo", PluginOutput: "OK", Timestamp: 1484527962, State: 0}
	if err := queueData(p); err != nil {
		t.Errorf("queueData returned: %s", err)
	}
	if l := q.Len(); l != 1 {
		t.Errorf("Queue expecting to contain 1 element. Contains %d", l)
	}
}

func TestCacheWorker(t *testing.T) {
	initCache()
	q = lfc.NewQueue()
	if q.Len() > 0 {
		t.Error("Queue not empty on start")
	}

	testCases := []struct {
		host      string
		service   string
		output    string
		timestamp uint32
		state     int16
	}{
		{"host01", "service foo", "Bar", 1484527962, 0},
		// Update the same service entry as created previously
		{"host01", "service foo", "New output", 1484527963, 1},
		// Second service on the same host
		{"host01", "service bar", "OK", 1484527964, 0},
		// Another host
		{"host02", "service whatever", "OK", 1484527965, 0},
		// 2nd service with the same name as the other host
		{"host02", "service bar", "OK", 1484527966, 0},
	}

	for _, tt := range testCases {
		p := &nsca.DataPacket{HostName: tt.host, Service: tt.service, PluginOutput: tt.output, Timestamp: tt.timestamp, State: tt.state}
		q.Enqueue(p)
		cacheWorker(false)
		h, ok := cache[tt.host]
		if !ok {
			t.Errorf("Entry %s should be present", tt.host)
		}
		s, ok := h[tt.service]
		if !ok {
			t.Errorf("%s should still have a %s entry", tt.host, tt.service)
		}
		if s.timestamp != tt.timestamp {
			t.Errorf("entry '%s' has wrong timestamp. Got %d, expecting %d", tt.service, s.timestamp, tt.timestamp)
		}
		if s.output != tt.output {
			t.Errorf("entry '%s' has wrong output. Got %s, expecting %s", tt.service, s.output, tt.output)
		}
		if s.state != tt.state {
			t.Errorf("entry '%s' has wrong state. Got %d, expecting %d", tt.service, s.state, tt.state)
		}
	}
}
