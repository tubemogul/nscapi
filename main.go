package main

import (
	"flag"
	"github.com/PurpureGecko/go-lfc"
	nsca "github.com/tubemogul/nscatools"
	"os"
	"strconv"
	"time"
)

var (
	q   *lfc.Queue
	cfg struct {
		apiIP          string
		apiPort        uint
		nscaIP         string
		nscaPort       uint
		nscaPassword   string
		nscaEncryption uint
	}
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

// getStringFromEnv gets the string value of the specified environment variable
// or the default value if this variable is not set
func getStringFromEnv(varName string, defaultValue string) string {
	var (
		val     string
		present bool
	)
	if val, present = os.LookupEnv(varName); !present {
		val = defaultValue
	}
	return val
}

// getUintFromEnv gets the uint value of the specified environment variable
// or the default value if this variable is not set
func getUintFromEnv(varName string, defaultValue uint, bitSize int) uint {
	var (
		strVal  string
		val     uint
		present bool
	)
	strVal, present = os.LookupEnv(varName)
	if val64, err := strconv.ParseUint(strVal, 10, bitSize); !present || err != nil {
		val = defaultValue
	} else {
		val = uint(val64)
	}
	return val
}

// initConfig initializes the configuration based on environment variables and
// passed flags (in the 12 factor app spirit)
func initConfig() {
	flag.StringVar(&cfg.apiIP, "api-ip", getStringFromEnv("NSCAPI_API_IP", "0.0.0.0"), "IP the API should listen on. Default to the NSCAPI_API_IP environment variable. Fallback: 0.0.0.0")
	flag.UintVar(&cfg.apiPort, "api-port", getUintFromEnv("NSCAPI_API_PORT", 8080, 32), "Port the API should listen on. Default to the NSCAPI_API_PORT environment variable. Fallback: 8080")
	flag.StringVar(&cfg.nscaIP, "nsca-server-ip", getStringFromEnv("NSCAPI_NSCA_IP", "0.0.0.0"), "IP the NSCA server should listen on. Default to the NSCAPI_NSCA_IP environment variable. Fallback: 0.0.0.0")
	flag.UintVar(&cfg.nscaPort, "nsca-server-port", getUintFromEnv("NSCAPI_NSCA_PORT", 5667, 16), "Port the NSCA server should listen on. Default to the NSCAPI_NSCA_PORT environment variable. Fallback: 5667")
	flag.StringVar(&cfg.nscaPassword, "nsca-server-password", getStringFromEnv("NSCAPI_NSCA_PASSWORD", ""), "Password the NSCA server should use. Default to the NSCAPI_NSCA_PASSWORD environment variable. Fallback: ''")
	flag.UintVar(&cfg.nscaEncryption, "nsca-server-encryption", getUintFromEnv("NSCAPI_NSCA_ENCYPTION", 0, 8), "Number corresponding to the encryption to be used by the NSCA server. Default to the NSCAPI_NSCA_ENCRYPTION environment variable. Fallback: 0. See 'DECRYPTION METHOD' on https://github.com/NagiosEnterprises/nsca/blob/master/sample-config/nsca.cfg.in for more details. Must be <27.")
	flag.Parse()
}

func main() {
	q = lfc.NewQueue()
	initCache()

	// Loads config from flags or from env
	initConfig()

	// Start the worker that updates the cache
	go cacheWorker(true)

	// Start the API inside a routine
	go initAPIServer(cfg.apiIP, cfg.apiPort)

	// Start the nsca server
	cfg := nsca.NewConfig(cfg.nscaIP, uint16(cfg.nscaPort), int(cfg.nscaEncryption), cfg.nscaPassword, queueData)
	nsca.StartServer(cfg, false)

}
