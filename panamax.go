package godutch

import (
	"errors"
	"github.com/thejerf/suture"
	"log"
)

//
// A Panamax is a world-wide known container cargo vessel. The name is adopted
// here, a local group of structures to register Containers and inventory their
// checks, also interfacing the final check execution call.
//
type Panamax struct {
	*suture.Supervisor
	Containers   map[string]*Container
	Checks       map[string]string
	ServiceToken map[string]suture.ServiceToken
}

// Creates a new Panamax instance, loading a Supervisor instance and allocating
// memory for internal registers.
func NewPanamax() *Panamax {
	var p *Panamax
	p = &Panamax{
		Supervisor: suture.New("Panamax", suture.Spec{
			Log: func(line string) {
				log.Println("Suture:", line)
			},
		}),
		Containers:   make(map[string]*Container),
		Checks:       make(map[string]string),
		ServiceToken: make(map[string]suture.ServiceToken),
	}
	return p
}

// Loads a Container onto Panamax, which means registering on local structs,
// listing the check names belonging to the Container. This method will execute
// bootstrap if container is not yet on this state.
func (p *Panamax) Onboard(c *Container) error {
	var err error
	var checkName string

	// bootstrapping the container if it's not loaded yet
	if !c.Bootstrapped {
		if err = c.Bootstrap(); err != nil {
			return err
		}
	}

	// a container must have at least one check
	if len(c.Checks) <= 0 {
		err = errors.New("Container '" + c.Name + "' has no Checks.")
		log.Fatalln(err)
		return err
	}

	log.Println("Loading container:", c.Name)

	p.Containers[c.Name] = c
	for _, checkName = range c.Checks {
		log.Println("Loading check:", checkName)
		p.Checks[checkName] = c.Name
	}

	return nil
}

// Adding background process to the local Supervisor, saving the unique
// service-id into the local registry.
func (p *Panamax) RegisterService(c *Container) {
	p.ServiceToken[c.Name] = p.Add(c.Bg)
}

// Removes a Container from Panamax, where the only input is the container name.
// The regarded process will also be shutdown.
func (p *Panamax) Offboard(containerName string) error {
	var err error
	var c *Container
	var okay bool = false
	var checkName string

	if c, okay = p.Containers[containerName]; !okay {
		err = errors.New("Can't find container '" + containerName + "'")
		return err
	}

	// removing container background process from supervisor
	p.Remove(p.ServiceToken[containerName])

	// cleaning up container and checks
	delete(p.Containers, containerName)
	delete(p.ServiceToken, containerName)
	for _, checkName = range c.Checks {
		delete(p.Checks, checkName)
	}

	if err = c.Shutdown(); err != nil {
		log.Fatalln("Error on shutting down container:", err)
	}

	return nil
}

// Call for check execution on it's respective container, creating then a
// request carrying command and arguments from this method pameters.
func (p *Panamax) Execute(cmd string, args []string) (*Response, error) {
	var err error
	var c *Container
	var containerName string
	var req []byte
	var resp *Response
	var okay bool

	if containerName, okay = p.Checks[cmd]; okay {
		log.Println("Container:", containerName, ", Command:", cmd)
		c = p.Containers[containerName]
	} else {
		err = errors.New("Can't find check '" + cmd + "' on any container.")
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
