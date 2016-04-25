package godutch

//
// GoDutch is the primary type on this software, acting on the integration
// position of network and local containers.
//

import (
	gocache "github.com/patrickmn/go-cache"
	"log"
	"time"
)

//
// Holds the references of configuration, Panamax and NRPE service, linking
// those elements to work together.
//
type GoDutch struct {
	cfg   *Config
	p     *Panamax
	cache *gocache.Cache
	ns    *NrpeService
	cs    *CarbonService
}

// Instantiates a new GoDutch, which will also spawn a new Panamax.
func NewGoDutch(cfg *Config) (*GoDutch, error) {
	var cache *gocache.Cache
	var p *Panamax
	var err error

	cache = gocache.New(time.Minute, 20*time.Second)

	if p, err = NewPanamax(cache); err != nil {
		return nil, err
	}

	return &GoDutch{cfg: cfg, p: p, cache: cache, ns: nil}, nil
}

// Go through the configured containers and load (unless disabled).
func (g *GoDutch) LoadContainers() error {
	var name string
	var containerCfg *ContainerConfig
	var err error

	for name, containerCfg = range g.cfg.Container {
		log.Printf("[GoDutch] Container: '%s'", name)
		if !containerCfg.Enabled {
			log.Printf("[GoDutch] Skipping, Container is disabled.")
			continue
		}
		if err = g.p.Load(containerCfg); err != nil {
			log.Printf("[GoDutch] Error loading container '%s'", name)
			return err
		}
	}

	return nil
}

// Loads all services listed on configuration files, skips when it's disabled
// and had specific loading mechanisms for each service. Return error.
func (g *GoDutch) LoadServices() error {
	var serviceCfg *ServiceConfig
	var name string
	// var err error

	for name, serviceCfg = range g.cfg.Service {
		log.Printf("[GoDutch] Service: '%s' (%s)", name, serviceCfg.Type)

		if !serviceCfg.Enabled {
			log.Printf("[GoDutch] Skipping '%s' (%s), Service is disabled.",
				name, serviceCfg.Type)
			continue
		}

		switch serviceCfg.Type {
		case "nrpe":
			log.Println("[GoDutch] Loading NRPE Service")
			// initializing NRPE service and informing local Panamax instance,
			// then the service is able to call for checks execution
			g.ns = NewNrpeService(serviceCfg, g.p)
		case "nsca":
			log.Println("[GoDutch] Loading NSCA Service")
		case "carbon":
			log.Println("[GoDutch] Loading Carbon Relay Service")
			// spawning a new Carbon Relay type of service, using local cache to
			// dispatch metrics
			g.cs = NewCarbonService(serviceCfg, g.cache)

			// background routine to pick up items from cache and send to carbon
			go func(cs *CarbonService) {
				for {
					time.Sleep(5 * time.Second)
					g.cs.Send()
				}
			}(g.cs)
		default:
			panic("[GoDutch] Service type is unkown: " + serviceCfg.Type)
		}
	}

	return nil
}

// Wraps serve method on NRPE service.
func (g *GoDutch) Serve() {
	if g.ns == nil {
		panic("NRPE Service is not loaded, nothing to Serve.")
	}
	// nrpe service in background
	go g.ns.Serve()
}

// Wraps stop call for the NRPE service and Panamax objects.
func (g *GoDutch) Stop() {
	g.ns.Stop()
	g.p.Stop()
}

/* EOF */
