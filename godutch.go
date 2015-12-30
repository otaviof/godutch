package godutch

import (
	"errors"
	"github.com/thejerf/suture"
	"log"
)

//
// Interface to link different components of GoDutch.
//
type Composer interface {
	Bootstrap() error
	Shutdown() error
	Execute(req []byte) (*Response, error)
	ComponentInfo() *Component
}

type Component struct {
	Name     string
	Checks   []string
	Type     string
	Instance suture.Service
}

//
// Primary integration point to join to inventory Containers, Service and make
// it behave as a Supervisor (with Suture's help).
//
type GoDutch struct {
	*suture.Supervisor
	Containers map[string]Composer
	Checks     map[string]string
	Services   map[string]Composer
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
			// how resilient this supervisor should be
			FailureDecay:     1,
			FailureThreshold: 11,
			FailureBackoff:   1,
			Timeout:          3,
		}),
		Containers: make(map[string]Composer),
		Checks:     make(map[string]string),
		Services:   make(map[string]Composer),
		tokens:     make(map[string]suture.ServiceToken),
	}
	return g
}

// Identify and register the Container's checks into GoDutch.
func (g *GoDutch) onboardContainer(c Composer) error {
	var err error
	var checkName string
	var component *Component = c.ComponentInfo()

	// chance for loading it's process
	if err = c.Bootstrap(); err != nil {
		return err
	}

	log.Println("Debug -> checks:", component.Checks)
	// a container must have at least one check
	if len(component.Checks) <= 0 {
		err = errors.New("Container '" + component.Name + "' has no Checks.")
		return err
	}

	log.Println("Loading container:", component.Name)
	g.Containers[component.Name] = c
	for _, checkName = range component.Checks {
		log.Printf("** Loading check: '%s' (%s)", checkName, component.Name)
		g.Checks[checkName] = component.Name
	}

	return nil
}

// Bootstrap a service, the only step required to onboard a service type.
func (g *GoDutch) onboardService(c Composer) error {
	var err error
	var component *Component = c.ComponentInfo()

	g.Services[component.Name] = c
	if err = c.Bootstrap(); err != nil {
		return err
	}
	return nil
}

// Loads a Object implementing GoDutch interface, for Containers, it will load
// the actual checks, for a service, it will only keep references for Suture.
func (g *GoDutch) Onboard(c Composer) error {
	var err error
	var component *Component = c.ComponentInfo()

	switch component.Type {
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
		err = errors.New("Component type is not known: " + component.Type)
		return err
	}

	return nil
}

// Adding background process to the local Supervisor, saving the unique
// service-id into the local registry.
func (g *GoDutch) Register(c Composer) {
	var component *Component = c.ComponentInfo()
	g.tokens[component.Name] = g.Add(component.Instance)
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
	}

	// or it might be container then
	if err = g.offboardContainer(name); err != nil {
		log.Fatalln("Error on offboarding a container:", err)
		return err
	}

	// both have tokens, removing here
	delete(g.tokens, name)

	// removing container background process from supervisor
	g.Remove(g.tokens[name])

	return nil
}

// Does the steps to offload a Service
func (g *GoDutch) offboardService(name string) error {
	// presence is already check by Offboard method
	delete(g.Services, name)
	return nil
}

// Removes a Container from GoDutch, where the only input is the container name.
// The regarded process will also be shutdown.
func (g *GoDutch) offboardContainer(name string) error {
	var err error
	var c Composer
	var okay bool = false
	var checkName string
	var component *Component

	// loading check list from method
	if c, okay = g.Containers[name]; okay {
		component = c.ComponentInfo()
		delete(g.Containers, component.Name)
		for _, checkName = range component.Checks {
			delete(g.Checks, checkName)
		}
		if err = c.Shutdown(); err != nil {
			log.Fatalln("Error on shutting down container:", err)
		}
	} else {
		err = errors.New("Can't find container '" + component.Name + "'")
		return err
	}

	return nil
}

// Call for check execution on it's respective container, creating then a
// request carrying command and arguments from this method pameters.
func (g *GoDutch) Execute(cmd string, args []string) (*Response, error) {
	var err error
	var c Composer
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
		log.Println("On request:", string(req[:]))
		log.Println("Error on Execute '", cmd, "':", err)
		return nil, err
	}

	return resp, err
}

/* EOF */
