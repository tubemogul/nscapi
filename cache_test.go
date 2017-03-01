package main

import "testing"

func TestInitCache(t *testing.T) {
	if cache != nil {
		t.Errorf("Already initialized cache found")
	}
	initCache()
	if cache == nil {
		t.Errorf("Cache object still not initialized after the init")
	}
}

func TestUpdateCacheEntry(t *testing.T) {
	initCache()

	// Create new entry
	h, ok := cache["host01"]
	if ok {
		t.Errorf("Entry host01 already exists")
	}
	updateCacheEntry("host01", "service foo", "OK", 1484527962, 0)
	h, ok = cache["host01"]
	if !ok {
		t.Errorf("Entry host01 should have been created")
	}
	s, ok := h["service foo"]
	if !ok {
		t.Errorf("host01 should have a service foo entry")
	}
	if s.timestamp != 1484527962 {
		t.Errorf("service foo has wrong timestamp upon creation")
	}
	if s.statusFirstSeen != 1484527962 {
		t.Errorf("service foo has wrong statusFirstSeen upon creation")
	}
	if s.output != "OK" {
		t.Errorf("service foo has wrong output upon creation")
	}
	if s.state != 0 {
		t.Errorf("service foo has wrong state upon creation")
	}

	testData := []struct {
		host            string
		service         string
		output          string
		timestamp       uint32
		statusFirstSeen uint32
		state           int16
	}{
		// Update the same service entry as created previously
		{"host01", "service foo", "New output", 1484527963, 1484527963, 1},
		// Second service on the same host
		{"host01", "service bar", "OK", 1484527964, 1484527964, 0},
		// Another host
		{"host02", "service whatever", "OK", 1484527965, 1484527965, 0},
		// 2nd service with the same name as the other host
		{"host02", "service bar", "OK", 1484527966, 1484527966, 0},
	}
	// Launch all the updates before running the tests
	for _, tt := range testData {
		updateCacheEntry(tt.host, tt.service, tt.output, tt.timestamp, tt.state)
	}

	for _, tt := range testData {
		h, ok = cache[tt.host]
		if !ok {
			t.Errorf("Entry %s should be present", tt.host)
		}
		s, ok = h[tt.service]
		if !ok {
			t.Errorf("%s should still have a %s entry", tt.host, tt.service)
		}
		if s.statusFirstSeen != tt.statusFirstSeen {
			t.Errorf("entry '%s' has wrong statusFirstSeen. Got %d, expecting %d", tt.service, s.statusFirstSeen, tt.statusFirstSeen)
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

	// Check that the statusFirstSeen is kept when updating with the same state
	updateCacheEntry("host02", "service bar", "OK", 1488527969, 0)
	if s.statusFirstSeen != 1484527966 {
		t.Errorf("entry '%s' has wrong statusFirstSeen upon update with no state change. Got %d, expecting %d", "service bar", s.statusFirstSeen, 1484527966)
	}
}
