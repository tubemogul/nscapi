package main

import (
	"testing"
)

func TestStatusString(t *testing.T) {
	cases := []struct {
		in  int16
		out string
	}{{0, "OK"}, {1, "Warning"}, {2, "Critical"}, {3, "Unknown"}}
	for _, tt := range cases {
		returned := statusString(tt.in)
		if returned != tt.out {
			t.Errorf("statusString(%d) should return %s, not %s", tt.in, tt.out, returned)
		}
	}
}
