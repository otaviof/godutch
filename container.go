package godutch

import (
	"errors"
	"log"
	"net"
	"os/exec"
	"strings"
)

/** TODOs *********************************************************************
 * - Find a way to map a go-channel into the socket, the "Container" needs to
 *   watch if a given socket is alive;
 * - The Container knows which checks are available through
 *   "__list_check_methods" call;
 * - Think on how to wrap the protocol, the json mappings and conversion into
 *   []byte, which will be written into the socket;
 */

//
// A Container is a wrapper of tools around a background process, on which we
// communicate using a socket and GoDutch-Protocol, based on JSON.
//
type Container struct {
	Name   string
	Bg     *BgCmd
	socket net.Conn
}

func NewContainer(name string, command []string) (*Container, error) {
	var bg *BgCmd
	var c *Container

	if len(command) < 2 {
		err := errors.New("Informed command is not long enough:" +
			strings.Join(command, " "))
		return nil, err
	}

	log.Println("*** Container:", name, "***")
	log.Println("Container command:", strings.Join(command, " "))
	// expanding the command argument into BgCmd
	bg = NewBgCmd(name, exec.Command(command[0], command[1:]...))
	c = &Container{Name: name, Bg: bg}

	return c, nil
}

func (c *Container) Bootstrap() error {
	var err error

	log.Println("Bootstraping Container:", c.Name)
	log.Println("Container's socket path:", c.Bg.SocketPath)

	// creating a reader on background command's socket
	if c.socket, err = net.Dial("unix", c.Bg.SocketPath); err != nil {
		log.Fatalln("Dialing to socket error: '", err, "'")
		return err
	}

	if err = c.listCheckMethods(); err != nil {
		return err
	}

	return nil
}

func (c *Container) listCheckMethods() error {
	var req []byte
	var err error

	req, _ = NewRequest("__list_check_methods", []string{})
	if _, err = c.socket.Write(req); err != nil {
		log.Fatalln("Socket write error:", err)
		return err
	}
	return nil
}

/* EOF */
