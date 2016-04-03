package godutch

//
// GoDutch is the primary type on this software, acting on the integration
// position of network and local containers.
//

import (
	"log"
)

//
// Holds the references of configuration, Panamax and NRPE service, linking
// those elements to work together.
//
type GoDutch struct {
	cfg *Config
	p *Panamax
	ns *NrpeService
}

// Instantiates a new GoDutch, which will also spawn a new Panamax.
func NewGoDutch(cfg *Config) (*GoDutch, error) {
	var g *GoDutch
	var p *Panamax
	var err error

	if p, err = NewPanamax(); err != nil {
		return nil, err
	}

	g = &GoDutch{
		cfg: cfg,
		p: p,
	}

	return g, nil
}

// Go through the configured containers and load (unless disabled).
func (g *GoDutch) LoadContainers() (error) {
	var name string
	var containerCfg *ContainerConfig
	var err error

	for name, containerCfg = range g.cfg.Containers {
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

// Based on configuration loads the k
func (g *GoDutch) LoadNrpeService() (error) {
	var name string
	var serviceCfg *ServiceConfig
	// var err error

	for name, serviceCfg = range g.cfg.Services {
		log.Printf("[GoDutch] Service: '%s'", name)

		if serviceCfg.Type != "nrpe" {
			log.Printf("[GoDutch] Skipping, not 'nrpe' type of service.")
			continue
		}

		if !serviceCfg.Enabled {
			log.Printf("[GoDutch] Skipping, Service is disabled.")
			continue
		}

		// initializing NRPE service and informing local Panamax instance
		g.ns = NewNrpeService(serviceCfg, g.p)
		// only a single nrpe service instance will be loaded
		break
	}

	return nil
}

// Start to serve, as in listening on the network interface.
func (g *GoDutch) Serve() {
	g.ns.Serve()
}

/* EOF */
