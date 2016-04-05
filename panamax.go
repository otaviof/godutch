package godutch

//
// Panamax is a world class container carrier vessel, and here the analogy
// applies for holding the check containers using Supervisor.
// On GoDutch this has one more role, to route the check requests towards the
// right container, since it keeps the inventory.
//

import (
	"errors"
	gocache "github.com/patrickmn/go-cache"
	"github.com/thejerf/suture"
	"log"
	"time"
)

//
// Containers and Checks inventory, plus Supervisor structure.
//
type Panamax struct {
	*suture.Supervisor
	containers map[string]*Container
	checks     map[string]*Container
	cache      *gocache.Cache
}

// Creates a new Panamax instnace. Alocates memotry and loads a new supervisor
// instance to hold the Containers.
func NewPanamax(cache *gocache.Cache) (*Panamax, error) {
	var p *Panamax = &Panamax{
		Supervisor: suture.New("Panamax", suture.Spec{
			Log: func(line string) { log.Println("[SUTURE-Panamax]", line) },
		}),
		containers: make(map[string]*Container),
		checks:     make(map[string]*Container),
		cache:      cache,
	}

	// letting the Supervisor run in background right from the start, it will be
	// requested to be on running state when onboarding containers
	go p.ServeBackground()

	return p, nil
}

// Loads a container based on configuration, starting command in background and
// loading it's inventory right after. When Container has no checks it will
// return error.
func (p *Panamax) Load(cfg *ContainerConfig) error {
	var found bool = false
	var item string
	var err error

	log.Printf("[Panamax] Loading container: '%s'", cfg.Name)
	if _, found = p.containers[cfg.Name]; found {
		return errors.New("[Panamax] Container already loaded: " + cfg.Name)
	}

	if p.containers[cfg.Name], err = NewContainer(cfg); err != nil {
		return err
	}

	// loading container on local Supervisor and quick sleep, to give it time to
	// start and be able to respond
	p.Add(p.containers[cfg.Name].Client())
	time.Sleep(1e9)

	if err = p.containers[cfg.Name].Bootstrap(); err != nil {
		log.Printf("[Panamax] Error on boostrapping container")
		return err
	}

	// having no checks found on this continer will return error
	if len(p.containers[cfg.Name].Inventory()) <= 0 {
		err = errors.New("[Panamax] No inventory found on: " + cfg.Name)
		return err
	}

	// loading container inventory
	for _, item = range p.containers[cfg.Name].Inventory() {
		log.Printf("[Panamax] Container '%s' has check: '%s'", cfg.Name, item)
		p.checks[item] = p.containers[cfg.Name]
	}

	return nil
}

// Wraps the Execute method from the Container using local inventory, save the
// results into Cache.
func (p *Panamax) Execute(req *Request) (*Response, error) {
	var name string = req.Fields.Command
	var found bool = false
	var resp *Response
	var err error

	// check's command is it's name, can be found on Request's fields
	if _, found = p.checks[name]; !found {
		log.Printf("[Panamax] Can't find check named '%s'", name)
		err = errors.New("[Panamax] Can't find a check named:" + name)
		return nil, err
	}

	if resp, err = p.checks[name].Execute(req); err != nil {
		return nil, err
	}

	// saving object on cache
	p.cache.Set(name, resp, gocache.DefaultExpiration)
	log.Printf("[Panamax] Cache count: '%d'", p.cache.ItemCount())

	return resp, nil
}

/* EOF */
