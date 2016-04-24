package godutch

//
// Implements a type of service that reads from local cache and offload this
// information into a Carbon relay server.
//

import (
	"fmt"
	gocarbon "github.com/jforman/carbon-golang"
	gocache "github.com/patrickmn/go-cache"
	"log"
)

type CarbonService struct {
	cfg    *ServiceConfig
	cache  *gocache.Cache
	DialOn []string
}

// Creates a new instance of CarbonService, which takes a cache object.
func NewCarbonService(cfg *ServiceConfig, cache *gocache.Cache) (
	*CarbonService, error) {
	var cs *CarbonService
	cs = &CarbonService{
		cfg:    cfg,
		cache:  cache,
		DialOn: cfg.ParseDialOn(),
	}
	return cs, nil
}

// Sends the metrics towards carbon server, first ask for gathering of the
// values that will be transferred. It tries on the configured server end-points
// sequentially, logging the results.
func (cs *CarbonService) Send() error {
	var err error
	var host string
	var port int
	var dialStr string
	var carbon *gocarbon.Carbon
	var i int
	var last int = len(cs.DialOn)
	var metrics []gocarbon.Metric

	metrics = cs.extractMetricsFromCache()

	for i, dialStr = range cs.DialOn {
		// extracting host and port from the dial-string
		host, port = cs.cfg.ParseDialString(dialStr)
		log.Printf("[Carbon] Connecting to: '%s:%d'", host, port)

		// instantiating Carbon, which will try to connect immediately, and
		// hence we can capture here connection errors, and try to use another
		// host on the list
		if carbon, err = gocarbon.NewCarbon(host, port, false, false); err != nil {
			log.Printf("[Carbon] Error on connecting to: '%s:%d'", host, port)
			log.Println("[Carbon] Error:", err)

			// checking if there's more hosts to try connection
			if i >= last {
				log.Println("[Carbon] No more hosts to try.")
				return err
			} else {
				continue
			}
		}

		log.Printf("[Carbon] Sending '%d' metric(s) towards '%s:%d'",
			len(metrics), host, port)

		if err = carbon.SendMetrics(metrics); err != nil {
			log.Println("[Carbon] Send metrics returned error:", err)
			continue
		} else {
			log.Println("[Carbon] Metrics sent!")
			break
		}
	}

	return nil
}

// Search for cached items and their respective metrics to be sent into Carbon
// service, cache object can't be expired and shall contain metrics before being
// picked up.
func (cs *CarbonService) extractMetricsFromCache() []gocarbon.Metric {
	var itemName string
	var item gocache.Item
	var cached interface{}
	var found bool
	var resp *Response
	var metric map[string]int
	var metricName string
	var metricValue int
	var metrics []gocarbon.Metric

	for itemName, item = range cs.cache.Items() {
		log.Printf("[Carbon] Reading from cache: '%s'", itemName)

		if item.Expired() {
			log.Printf("[Carbon] Cache item is expired: '%s'", itemName)
			continue
		}

		// loading Response object from Cache
		if cached, found = cs.cache.Get(itemName); !found {
			log.Printf("[Carbon] Key is not found on Cache: '%s'", itemName)
			continue
		} else {
			// transforming from interface back into Response type
			resp = cached.(*Response)
		}

		// checking whether are metrics to be sent
		if len(resp.Metrics) <= 0 {
			log.Printf("[Carbon] Cache entry has no metrics.")
			continue
		}

		// finally, collecting the metrics
		for _, metric = range resp.Metrics {
			for metricName, metricValue = range metric {
				log.Printf("[Carbon] Collecting metric: '%s.%s' -> %d",
					itemName, metricName, metricValue)

				metrics = append(
					metrics,
					gocarbon.Metric{
						Name:      fmt.Sprintf("%s.%s", itemName, metricName),
						Value:     float64(metricValue),
						Timestamp: int64(resp.Ts),
					},
				)
			}
		}
	}

	return metrics
}

/* EOF */