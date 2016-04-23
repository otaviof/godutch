package godutch

//
// This is the local network listener for NRPE service, which has a integration
// with Panamax object, therefore it's able to handle the requests directly.
//

import (
	"fmt"
	"log"
	"net"
)

//
// NRPE service type, basically holds configuration.
//
type NrpeService struct {
	listener net.Listener
	cfg      *ServiceConfig
	p        *Panamax
	listenOn string
}

// Creates a new instance of NRPE serice, which recieves a pointer of Panamax,
// and then it's able to call checks on running containers.
func NewNrpeService(cfg *ServiceConfig, p *Panamax) *NrpeService {
	var ns *NrpeService
	ns = &NrpeService{
		cfg:      cfg,
		p:        p,
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
			log.Println("[Nrpe] Error on accepting connection:", err)
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
	if resp, err = ns.panamaxExecute(cmd, args); err != nil {
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
// towards Panamax.
func (ns *NrpeService) panamaxExecute(cmd string, args []string) (*Response, error) {
	var req *Request
	var resp *Response
	var err error

	if req, err = NewRequest(cmd, args); err != nil {
		return nil, err
	}

	if resp, err = ns.p.Execute(req); err != nil {
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
