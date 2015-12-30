package godutch

import (
	"bytes"
	"errors"
	"github.com/thejerf/suture"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

//
// A Container is a wrapper of tools around a background process, on which we
// communicate using a socket and GoDutch-Protocol, based on JSON.
//
type Container struct {
	Name         string
	Bg           *BgCmd
	socket       net.Conn
	bootstrapped bool
	Checks       []string
}

// Creates a new Container type, using name and command informed by parameter
// to spawn a background process via BgCmd, exposed using Bg attribute.
func NewContainer(name string, command []string) (*Container, error) {
	var bg *BgCmd

	if len(command) < 2 {
		err := errors.New("Informed command is not long enough:" +
			strings.Join(command, " "))
		return nil, err
	}

	log.Println("*** Container:", name, "***")
	log.Println("Container command:", strings.Join(command, " "))

	// expanding the command argument into BgCmd
	bg = NewBgCmd(name, command)

	return &Container{Name: name, Bg: bg}, nil
}

// Returns the component type, which in this case is a container
func (c *Container) ComponentType() string {
	return "container"
}

// Return the container name to "component" interface.
func (c *Container) ComponentName() string {
	return c.Name
}

// Returns the object that shall be kept alive by Supervisor Trees (Suture), in
// the Container case we are interesed on having bgcmd object alive.
func (c *Container) ComponentObject() suture.Service {
	return c.Bg
}

// Part of "component" interface, list which checks are available.
func (c *Container) ComponentChecks() []string {
	return c.Checks
}

// Prepare a container to be up and running, opening the socket using
// Container's "socket" attribute.
func (c *Container) Bootstrap() error {
	var err error

	if c.bootstrapped {
		log.Println("Container has already been bootstraped.")
		return nil
	}

	log.Println("Bootstraping Container:", c.Name)
	log.Println("Container's socket path:", c.Bg.SocketPath)

	if err = c.listCheckMethods(); err != nil {
		return err
	}

	c.bootstrapped = true

	return nil
}

func (c *Container) socketDial() error {
	var err error
	var counter int = 0

	for {
		counter += 1
		// creating a reader on background command's socket
		if c.socket, err = net.Dial("unix", c.Bg.SocketPath); err != nil {
			log.Println("(", counter, "/ 3 ) net.Dial error: '", err, "'")
			if counter >= 3 {
				return err
			} else {
				time.Sleep(time.Second)
				continue
			}
		}
		return nil
	}
}

// Stop a container, closing the socket and asking os/exec to kill the process,
// if not dead just yet.
func (c *Container) Shutdown() error {
	defer c.socket.Close()
	c.Bg.Stop()
	return nil
}

// Executes the "__list_check_methods" call on the socket interface, load the
// response onto Container's Checks array of strings.
func (c *Container) listCheckMethods() error {
	var req []byte
	var resp *Response
	var err error

	req, _ = NewRequest("__list_check_methods", []string{})
	if resp, err = c.Execute(req); err != nil {
		log.Fatalln("Socket write error:", err)
		return err
	}

	log.Println("Available checks:", strings.Join(resp.Stdout, ", "))
	c.Checks = resp.Stdout

	return nil
}

// Execute a request towards the socket interface, simple by syncronously
// writing on the socket, and via a goroutine reading back from it, which must
// be a Response type of payload.
func (c *Container) Execute(req []byte) (*Response, error) {
	var err error
	var payload []byte
	var resp *Response
	var respCh chan []byte = make(chan []byte)
	var errorCh chan error = make(chan error)

	if c.socketDial(); err != nil {
		log.Fatalln("Socket DIAL error:", err)
		return nil, err
	}

	log.Println("Sending request:", string(req[:]))
	if _, err = c.socket.Write(req); err != nil {
		log.Println("Socket WRITE error:", err)
		return nil, err
	}

	// background routine to read socke's FD, informing response and error
	// channels when there's data back
	go c.socketReader(respCh, errorCh)

	defer c.socket.Close()

	// TODO
	//   * Handle request timeouts;
	// http://stackoverflow.com/questions/9680812
	for {
		select {
		case payload = <-respCh:
			log.Println("Got back:", string(payload[:]))
			// converting JSON contents into local struct
			if resp, err = NewResponse(payload[:]); err != nil {
				log.Fatalln("Parsing response:", err)
				return nil, err
			}
			log.Printf("Response: %#v", resp)
			return resp, nil
		case err = <-errorCh:
			log.Println("Socket reading error:", err)
			return nil, err
		}
	}
}

// Reads from a socket file descriptor onto a local buffer, which is by the end
// sent to response-channel (respCh), informed by parameters. Error is captured
// locally and also sent back by error-channel (errorCh).
func (c *Container) socketReader(respCh chan []byte, errorCh chan error) {
	var buf bytes.Buffer
	for {
		_, err := io.Copy(&buf, c.socket)
		if err != nil {
			log.Println("Socket read error:", err)
			errorCh <- err
			return
		}
		respCh <- buf.Bytes()
	}
}

/* EOF */
