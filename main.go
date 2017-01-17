package main

import (
	"github.com/PurpureGecko/go-lfc"
	nsca "github.com/tubemogul/nscatools"
	"log"
	"os"
	"time"
)

var (
	dbg *log.Logger
	q   *lfc.Queue
)

// cacheWorker will continuously pull DataPackets out of the queue and update
// the cache with it. Passing true as parameter will make it loop indefinetly.
// Passing false as parameter will make it return whenever the queue is empty
func cacheWorker(runIndefinetly bool) {
	for {
		if q.Len() == 0 {
			if runIndefinetly {
				// Wait 100ms before checking for data in the queue
				time.Sleep(100 * time.Millisecond)
			} else {
				return
			}
		} else {
			pkt, ok := q.Dequeue()
			if ok {
				if p, ok := pkt.(*nsca.DataPacket); ok {
					updateCacheEntry(p.HostName, p.Service, p.PluginOutput, p.Timestamp, p.State)
				}
			}
		}
	}
}

// queueData will put the DataPacket received by the nsca server in a
// non-locking queue
func queueData(p *nsca.DataPacket) error {
	q.Enqueue(p)
	return nil
}

func main() {
	debugHandle := os.Stdout
	dbg = log.New(debugHandle, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)

	q = lfc.NewQueue()
	initCache()

	// Loads config from file or from env

	// Start the worker that updates the cache
	go cacheWorker(true)

	// Start the API inside a routine
	go initAPIServer("localhost", 8080)

	// Start the nsca server
	cfg := nsca.NewConfig("localhost", 5667, nsca.Encrypt3DES, "toto", queueData)
	nsca.StartServer(cfg, true)

}
