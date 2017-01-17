package main

// The cache structure contains 2 layers of maps:
// * Layer 1: the key is the hostname
// * Layer 2: for each hostname, there's a map where the key is the service name
var cache map[string]map[string]*serviceEntry

// ServiceEntry can be found in the 2nd layer of he map and contains the details
// of the last status of the check (timestamp, state and plugin output)
type serviceEntry struct {
	timestamp uint32
	state     int16
	output    string
}

// initCache initialize the cache object
func initCache() {
	cache = make(map[string]map[string]*serviceEntry)
}

// updateCacheEntry adds or update a given service check result in the cache map
func updateCacheEntry(hostname, servicename, output string, timestamp uint32, state int16) {
	svc, ok := cache[hostname]
	if !ok {
		svc = make(map[string]*serviceEntry)
		cache[hostname] = svc
	}
	svc[servicename] = &serviceEntry{
		timestamp: timestamp,
		output:    output,
		state:     state,
	}
}
