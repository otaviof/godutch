package godutch

import (
	"errors"
	"github.com/thejerf/suture"
	"log"
)

type ContainerShip struct {
	*suture.Supervisor
	Containers map[string]*Container
	Checks     map[string]string
}

func NewContainerShip() *ContainerShip {
	var cs *ContainerShip
	cs = &ContainerShip{
		Supervisor: suture.New("ContainerShip", suture.Spec{
			Log: func(line string) {
				log.Println("Suture:", line)
			},
		}),
		Containers: make(map[string]*Container),
		Checks:     make(map[string]string),
	}
	return cs
}

func (cs *ContainerShip) Onboard(container *Container) error {
	var err error

	if len(container.Checks) <= 0 {
		err = errors.New("Container '" + container.Name + "' has no Checks.")
		log.Fatalln(err)
		return err
	}

	log.Println("Loading container:", container.Name)
	cs.Containers[container.Name] = container
	for _, checkName := range container.Checks {
		log.Println("Loading check:", checkName)
		cs.Checks[checkName] = container.Name
	}

	return err
}

func (cs *ContainerShip) Offboard(containerName string) error {
	var err error
	return err
}

func (cs *ContainerShip) Execute(cmd string, args []string) (*Response, error) {
	var err error
	var containerName string
	var c *Container
	var resp *Response
	var req []byte

	containerName = cs.Checks[cmd]
	log.Println("Container:", containerName, ", Command:", cmd)
	c = cs.Containers[containerName]

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
