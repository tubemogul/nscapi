package main

// The cache structure contains 2 layers of maps:
// * Layer 1: the key is the hostname
// * Layer 2: for each hostname, there's a map where the key is the service name
var cache map[string]map[string]*serviceEntry

// ServiceEntry can be found in the 2nd layer of he map and contains the details
// of the last status of the check (timestamp, timestamp of the last status
// change, state and plugin output)
type serviceEntry struct {
	timestamp       uint32
	statusFirstSeen uint32
	state           int16
	output          string
}

// initCache initialize the cache object
func initCache() {
	cache = make(map[string]map[string]*serviceEntry)
}

// updateCacheEntry adds or update a given service check result in the cache map
func updateCacheEntry(hostname, servicename, output string, timestamp uint32, state int16) {
	svc, ok := cache[hostname]
	firstSeen := timestamp
	if !ok {
		svc = make(map[string]*serviceEntry)
		cache[hostname] = svc
	} else if service, exists := svc[servicename]; exists {
		// If the entry already exists and we update it, we want the time we've seen
		// the switch to the current status
		if service.state == state {
			firstSeen = service.statusFirstSeen
		}
	}
	svc[servicename] = &serviceEntry{
		timestamp:       timestamp,
		statusFirstSeen: firstSeen,
		output:          output,
		state:           state,
	}
}
