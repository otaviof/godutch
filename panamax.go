package godutch

import (
	"errors"
	"github.com/thejerf/suture"
	"log"
)

type Panamax struct {
	*suture.Supervisor
	Containers map[string]*Container
	Checks     map[string]string
}

func NewPanamax() *Panamax {
	var p *Panamax
	p = &Panamax{
		Supervisor: suture.New("Panamax", suture.Spec{
			Log: func(line string) {
				log.Println("Suture:", line)
			},
		}),
		Containers: make(map[string]*Container),
		Checks:     make(map[string]string),
	}
	return p
}

func (p *Panamax) Onboard(container *Container) error {
	var err error

	if len(container.Checks) <= 0 {
		err = errors.New("Container '" + container.Name + "' has no Checks.")
		log.Fatalln(err)
		return err
	}

	log.Println("Loading container:", container.Name)
	p.Containers[container.Name] = container
	for _, checkName := range container.Checks {
		log.Println("Loading check:", checkName)
		p.Checks[checkName] = container.Name
	}

	return err
}

func (p *Panamax) Offboard(containerName string) error {
	var err error
	return err
}

func (p *Panamax) Execute(cmd string, args []string) (*Response, error) {
	var err error
	var containerName string
	var c *Container
	var resp *Response
	var req []byte

	containerName = p.Checks[cmd]
	log.Println("Container:", containerName, ", Command:", cmd)
	c = p.Containers[containerName]

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
