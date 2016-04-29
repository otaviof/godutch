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
	// configuration object, for all assets
	cfg *Config
	// Panamax object, the container holder
	p *Panamax
	// cache instance, to dub as a shared object store
	cache *gocache.Cache
	// NRPE Service, listens to network and run checks
	ns *NrpeService
	// investigate cache for metrics and feed Carbon server
	cs *CarbonService
	// maximum threshold for running a check
	lastRunThreshold int64
}

// Instantiates a new GoDutch, which will also spawn a new Panamax.
func NewGoDutch(cfg *Config) (*GoDutch, error) {
	var cache *gocache.Cache
	var p *Panamax
	var g *GoDutch
	var err error

	cache = gocache.New(time.Minute, 20*time.Second)

	if p, err = NewPanamax(cache); err != nil {
		return nil, err
	}

	g = &GoDutch{
		cfg:              cfg,
		p:                p,
		cache:            cache,
		ns:               nil,
		cs:               nil,
		lastRunThreshold: -1,
	}

	return g, nil
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

	for name, serviceCfg = range g.cfg.Service {
		log.Printf("[GoDutch] Service: '%s' (%s)", name, serviceCfg.Type)

		if !serviceCfg.Enabled {
			log.Printf("[GoDutch] Skipping '%s' (%s), Service is disabled.",
				name, serviceCfg.Type)
			continue
		}

		log.Printf("[GoDutch] DEBUG: name: '%s', lastRunThreshold: '%d'",
			name, serviceCfg.LastRunThreshold)

		// when last-run-threshold is found, saving the lowest item
		if serviceCfg.LastRunThreshold > 0 {
			if g.lastRunThreshold == -1 || g.lastRunThreshold > serviceCfg.LastRunThreshold {
				log.Printf("[GoDutch] LastRunThreshold: from %ds to %ds'.",
					g.lastRunThreshold, serviceCfg.LastRunThreshold)
				g.lastRunThreshold = serviceCfg.LastRunThreshold
			}
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
		case "sensu":
			log.Println("[GoDutch] Loading Sensu Service")
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

	// nrpe service accepting connections
	go g.ns.Serve()

	// carbon relay inpecting cache and sending metrics
	go g.cs.Serve()

	// running check's that are delayed on shedule
	go g.runDelayedChecks()
}

// Wraps stop call for the NRPE service and Panamax objects.
func (g *GoDutch) Stop() {
	// nrpe service stop
	g.ns.Stop()
	// panamax (and it's containers) stop
	g.p.Stop()
}

func (g *GoDutch) runDelayedChecks() {
	var name string
	var lastRun int64
	var req *Request
	var err error

	if g.lastRunThreshold <= 0 {
		log.Println("[GoDutch] lastRunThreshold is disabled, no auto-run.")
		return
	}

	log.Printf("[GoDutch] Checks will run automatically after: %ds",
		g.lastRunThreshold)

	for {
		time.Sleep(10 * time.Second)
		log.Println("[GoDutch] Sleep done, looking for delayed checks.")

		for name, lastRun = range g.p.ChecksRunReport(g.lastRunThreshold) {
			log.Printf("[GoDutch] Executing '%s', last run at %ds ago (%ds threshold)",
				name, lastRun, g.lastRunThreshold)

			if req, err = NewRequest(name, []string{}); err != nil {
				log.Fatalln("[GoDutch] Error on creating request to: '%s'", name)
			}

			if _, err = g.p.Execute(req); err != nil {
				log.Println("[GoDutch] Error on execute: ", err)
			}
		}
	}
}

/* EOF */
