package godutch

import (
	"bytes"
	"errors"
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
	respCh       chan []byte
	errorCh      chan error
}

// Creates a new container with a background command.
func NewContainer(containerCfg *ContainerConfig) (*Container, error) {
	var err error
	var c *Container

	// verifying if socket directory exists
	if _, err = exists(containerCfg.SocketDir); err != nil {
		log.Fatalln(
			"[Container] Can't find socket directory: ('",
			containerCfg.SocketDir, "'):", err)
		return nil, err
	}

	if len(containerCfg.Command) < 2 {
		err = errors.New("Informed command is not long enough:" +
			strings.Join(containerCfg.Command, " "))
		return nil, err
	}

	log.Printf("[Container] Name: '%s', Command: '%s'",
		containerCfg.Name,
		strings.Join(containerCfg.Command, " "))

	c = &Container{
		Name: containerCfg.Name,
		Bg:   NewBgCmd(containerCfg),
	}

	return c, nil
}

// Method to retrieve informatoin about current component towards composer
// interface.
func (c *Container) ComponentInfo() *Component {
	return &Component{
		Name:     c.Name,
		Checks:   c.Checks,
		Type:     "container",
		Instance: c.Bg,
	}
}

// Prepare a container to be up and running, opening the socket using
// Container's "socket" attribute.
func (c *Container) Bootstrap() error {
	var err error

	if c.bootstrapped {
		log.Println("[Container] Already has been bootstraped, skipping.")
		return nil
	}

	log.Printf("[Container] Bootstraping: '%s', Socket path: '%s'",
		c.Name, c.Bg.SocketPath)

	if err = c.listCheckMethods(); err != nil {
		return err
	}

	c.bootstrapped = true

	return nil
}

// Dials to a socket using a counter to support a few attempts before just
// returning back the error.
func (c *Container) socketDial() error {
	var err error
	var counter int = 0

	for {
		counter += 1
		// creating a reader on background command's socket
		if c.socket, err = net.Dial("unix", c.Bg.SocketPath); err != nil {
			log.Println(
				"[Container] (", counter, "/ 3 ) net.Dial error: '", err, "'")
			// maximum retries before give up
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
		log.Fatalln("[Container] Socket write error:", err)
		return err
	}

	log.Printf("[Container] Checks: '%s'", strings.Join(resp.Stdout, "', '"))
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

	c.respCh = make(chan []byte)
	c.errorCh = make(chan error)

	if c.socketDial(); err != nil {
		log.Fatalln("[Container] Socket dial error:", err)
		return nil, err
	}

	log.Printf("[Container] Sending request: '%s'", string(req[:]))
	if _, err = c.socket.Write(req); err != nil {
		log.Println("[Container] Socket WRITE error:", err)
		return nil, err
	}

	// background routine to read socke's FD, informing response and error
	// channels when there's data back, for socket-close action we adopt defer
	go c.socketReader()
	// to be closed when we end this func, in other words, right after reading
	// data or handling connection error
	defer c.socket.Close()

	// TODO
	//  * Handle request timeouts (http://stackoverflow.com/questions/9680812);
	for {
		select {
		case payload = <-c.respCh:
			log.Printf("[Container] Request's payload: '%s'", string(payload[:]))
			if resp, err = NewResponse(payload[:]); err != nil {
				log.Fatalln("[Container] Error on parsing response:", err)
				return nil, err
			}
			return resp, nil
		case err = <-c.errorCh:
			log.Println("[Container] Socket reading error:", err)
			return nil, err
		}
	}
}

// Reads from a socket file descriptor onto a local buffer, which is by the end
// sent to response-channel (respCh), informed by parameters. Error is captured
// locally and also sent back by error-channel (errorCh).
func (c *Container) socketReader() {
	var err error
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, c.socket); err != nil {
		log.Println("[Container] Socket read error:", err)
		c.errorCh <- err
		return
	}
	c.respCh <- buf.Bytes()
}

/* EOF */
