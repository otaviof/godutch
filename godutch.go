package godutch

import (
	"errors"
	"github.com/thejerf/suture"
	"log"
)

//
// Interface to link different pieces of GoDutch. On which a Container will
// go though different onboarding steps than a Service.
//
type component interface {
	Bootstrap() error
	Shutdown() error
	Execute(req []byte) (*Response, error)
	// the "component" methods are the behavior modifier and the object that
	// will be kept alive in the background, using Suture
	ComponentName() string
	ComponentChecks() []string
	ComponentType() string
	ComponentObject() suture.Service
}

//
// A Panamax is a world-wide known container cargo vessel. The name is adopted
// here, a local group of structures to register Containers and inventory their
// checks, also interfacing the final check execution call.
//
type GoDutch struct {
	*suture.Supervisor
	Containers map[string]component
	Checks     map[string]string
	Services   map[string]component
	tokens     map[string]suture.ServiceToken
}

// Creates a new GoDutch instance, loading a Supervisor instance and allocating
// memory for internal registers.
func NewGoDutch() *GoDutch {
	var g *GoDutch
	g = &GoDutch{
		Supervisor: suture.New("GoDutch", suture.Spec{
			Log: func(line string) {
				log.Println("SUTURE:", line)
			},
		}),
		Containers: make(map[string]component),
		Checks:     make(map[string]string),
		Services:   make(map[string]component),
		tokens:     make(map[string]suture.ServiceToken),
	}
	return g
}

// Identify and register the Container's checks into GoDutch.
func (g *GoDutch) onboardContainer(c component) error {
	var err error
	var checkName string
	var name string
	var checks []string

	if err = c.Bootstrap(); err != nil {
		return err
	}

	// loading check list from method
	name = c.ComponentName()
	checks = c.ComponentChecks()

	// a container must have at least one check
	if len(checks) <= 0 {
		err = errors.New("Container '" + name + "' has no Checks.")
		return err
	}

	log.Println("Loading container:", name)

	g.Containers[name] = c
	for _, checkName = range checks {
		log.Printf("** Loading check: '%s' (%s)", checkName, name)
		g.Checks[checkName] = name
	}

	return nil
}

// Bootstrap a service, the only step required to onboard a service type.
func (g *GoDutch) onboardService(c component) error {
	var err error
	g.Services[c.ComponentType()] = c
	if err = c.Bootstrap(); err != nil {
		return err
	}
	return nil
}

// Loads a Object implementing GoDutch interface, for Containers, it will load
// the actual checks, for a service, it will only keep references for Suture.
func (g *GoDutch) Onboard(c component) error {
	var err error

	switch c.ComponentType() {
	case "container":
		if err = g.onboardContainer(c); err != nil {
			log.Fatalln("Errors on onboarding container:", err)
			return err
		}
	case "service":
		if err = g.onboardService(c); err != nil {
			log.Fatalln("Errors on onboarding service:", err)
			return err
		}
	default:
		err = errors.New("Component type is not known: " + c.ComponentType())
		return err
	}

	return nil
}

// Adding background process to the local Supervisor, saving the unique
// service-id into the local registry.
func (g *GoDutch) Register(c component) {
	g.tokens[c.ComponentName()] = g.Add(c.ComponentObject())
}

// Execute the Offboard of a Container or Service based on it's name.
func (g *GoDutch) Offboard(name string) error {
	var err error
	var okay bool = false

	// first let's check if this is a service
	if _, okay = g.Services[name]; okay {
		if err = g.offboardService(name); err != nil {
			log.Fatalln("Error on offboarding service:", err)
			return err
		}
		return nil
	} else {
		log.Printf("Offboard: '%s', not found on services!", name)
	}

	// or it might be container then
	if err = g.offboardContainer(name); err != nil {
		return err
	}

	// both have tokens, removing here
	delete(g.tokens, name)

	// removing container background process from supervisor
	g.Remove(g.tokens[name])

	return nil
}

// Does the steps to offload a Service, which implies on calling Shutdown and
// remove references from Services.
func (g *GoDutch) offboardService(name string) error {
	/*
		var err error
		if err = g.Services[name].Shutdown(); err != nil {
			log.Fatalln("Error on shutting down container:", err)
		}
	*/
	// presence is already check by Offboard method
	delete(g.Services, name)
	return nil
}

// Removes a Container from GoDutch, where the only input is the container name.
// The regarded process will also be shutdown.
func (g *GoDutch) offboardContainer(name string) error {
	var err error
	var c component
	var okay bool = false
	var checkName string
	var checks []string

	// loading check list from method
	if c, okay = g.Containers[name]; !okay {
		err = errors.New("Can't find container '" + name + "'")
		return err
	}

	checks = c.ComponentChecks()

	// cleaning up container and checks
	delete(g.Containers, name)

	for _, checkName = range checks {
		delete(g.Checks, checkName)
	}

	if err = c.Shutdown(); err != nil {
		log.Fatalln("Error on shutting down container:", err)
	}

	return nil
}

// Call for check execution on it's respective container, creating then a
// request carrying command and arguments from this method pameters.
func (g *GoDutch) Execute(cmd string, args []string) (*Response, error) {
	var err error
	var c component
	var containerName string
	var req []byte
	var resp *Response
	var okay bool

	log.Printf("GoDutch about to execute cmd: '%s'", cmd)

	if containerName, okay = g.Checks[cmd]; !okay {
		err = errors.New("Can't find check '" + cmd + "' on any container.")
		return nil, err
	}

	log.Println("Container:", containerName, ", Command:", cmd)
	if c, okay = g.Containers[containerName]; !okay {
		err = errors.New("Can't find container: " + containerName)
		return nil, err
	}

	if req, err = NewRequest(cmd, args); err != nil {
		log.Fatalln("Error on creating Request:", err)
		return nil, err
	}

	if resp, err = c.Execute(req); err != nil {
		log.Fatalln("On request:", string(req[:]))
		log.Fatalln("Error on Execute '", cmd, "':", err)
		return nil, err
	}

	return resp, err
}

/* EOF */
