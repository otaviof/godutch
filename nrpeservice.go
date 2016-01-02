package godutch

import (
	"fmt"
	"log"
	"net"
)

//
// Defines the NRPE service, which listens on configured interface and it's
// able to create a NrpePacket object with connection's payload.
//
type NrpeService struct {
	listener net.Listener
	cfg      *ServiceConfig
	g        *GoDutch
	listenOn string
}

// Creates a new object from NrpeService.
func NewNrpeService(cfg *ServiceConfig, g *GoDutch) *NrpeService {
	var ns *NrpeService
	ns = &NrpeService{
		cfg:      cfg,
		g:        g,
		listenOn: fmt.Sprintf("%s:%d", cfg.Interface, cfg.Port),
	}
	return ns
}

// Start listening on network interface and port, asyncronously will spawn a
// connection handler, when this event happen.
func (ns *NrpeService) Serve() {
	var err error
	var conn net.Conn

	log.Printf("[Nrpe] Listening on: '%s'", ns.listenOn)

	// creates a new network listener based on configuration
	if ns.listener, err = net.Listen("tcp", ns.listenOn); err != nil {
		log.Fatalln("[Nrpe] Error during net.Listen:", err)
		return
	}

	for {
		if conn, err = ns.listener.Accept(); err != nil {
			log.Println("[Nrpe] Error on accept connection:", err)
			return
		}
		go ns.handleConnection(conn)
	}
}

// Takes a network connection and extract it's buffer, passing along to create
// a NrpePacket, from which we can extract the actual command and it's
// arguments.
func (ns *NrpeService) handleConnection(conn net.Conn) {
	var err error
	var n int
	var buf []byte = make([]byte, NRPE_PACKET_SIZE)
	var pkt *NrpePacket
	var cmd string
	var args []string = []string{}
	var resp *Response

	if n, err = conn.Read(buf); n == 0 || err != nil {
		log.Println("[Nrpe] Error on reading from connection:", err)
		return
	}

	// transforming payload on a NRPE packet
	if pkt, err = NewNrpePacket(buf, n); err != nil {
		log.Println("[Nrpe] Payload:", string(buf[:]))
		log.Fatalln("[Nrpe] Error on response NRPE Packet:", err)
		return
	}

	if cmd, args, err = pkt.ExtractCmdAndArgsFromBuffer(); err != nil {
		log.Fatalln("[Nrpe] Error on parsing packet's buffer:", err)
		return
	}

	// using buffer to exectract command and it's argument
	if resp, err = ns.godutchExec(cmd, args); err != nil {
		log.Println("[Nrpe] Error on GODUTCH-EXEC:", err)
		resp = &Response{
			Name:   cmd,
			Status: STATE_UNKNOWN,
			Stdout: []string{fmt.Sprintf("[ERROR] %s", err)},
		}
	}

	// writing back to the connection
	if _, err = conn.Write(NrpePacketFromResponse(resp)); err != nil {
		log.Fatalln("[Nrpe] Error on writing response:", err)
		return
	}

	if err = conn.Close(); err != nil {
		log.Println("[Nrpe] Error on closing connection:", err)
	}
}

// Extract command and arguments from the packet buffer and compose a call
// towards GoDutch.
func (ns *NrpeService) godutchExec(cmd string, args []string) (*Response, error) {
	var err error
	var resp *Response
	if resp, err = ns.g.Execute(cmd, args); err != nil {
		return nil, err
	}
	return resp, nil
}

// Stop the service execution, which here for NRPE service means closing the
// network listener.
func (ns *NrpeService) Stop() {
	var err error
	if err = ns.listener.Close(); err != nil {
		log.Println("[Nrpe] Error on closing listener:", err)
	}
}

/* EOF */
