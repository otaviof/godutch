package godutch

import (
	"fmt"
	"log"
	"net"
)

type NRPESrvc struct {
	listener net.Listener
	cfg      *NRPEConfig
	p        *Panamax
}

func NewNRPESrvc(cfg *NRPEConfig, p *Panamax) (*NRPESrvc, error) {
	var err error
	var ns *NRPESrvc
	ns = &NRPESrvc{
		cfg: cfg,
		p:   p,
	}
	return ns, err
}

func (ns *NRPESrvc) Serve() {
	var err error
	var listenOn string = fmt.Sprintf("%s:%d", ns.cfg.Interface, ns.cfg.Port)
	var conn net.Conn

	// creates a new network listener based on configuration
	if ns.listener, err = net.Listen("tcp", listenOn); err != nil {
		log.Fatalln("Error during net.Listen:", err)
		return
	}

	for {
		if conn, err = ns.listener.Accept(); err != nil {
			log.Fatalln("Error during accept connection:", err)
			return
		}
		go ns.handleConnection(conn)
	}
}

func (ns *NRPESrvc) handleConnection(conn net.Conn) {
	var err error
	var n int
	var buf []byte = make([]byte, NRPE_PACKET_SIZE)
	var pkt *NRPEPacket
	var cmd string
	var args []string = []string{}
	var resp *Response

	if n, err = conn.Read(buf); n == 0 || err != nil {
		log.Fatalln("Error reading from connection:", err)
		return
	}
	defer conn.Close()
	log.Println("Bytes read from connection:", n)

	// transforming payload on a NRPE packet
	if pkt, err = NewNRPEPacket(buf, n); err != nil {
		log.Fatalln("Payload:", string(buf[:]))
		log.Fatalln("Error on instantiating a new NRPE Packet:", err)
		return
	}

	if cmd, args, err = pkt.ExtractCmdAndArgsFromBuffer(); err != nil {
		log.Fatalln("Error on extracting comamnd from buffer:", err)
		return
	}

	// using buffer to exectract command and it's argument
	if resp, err = ns.panamaxExec(cmd, args); err != nil {
		log.Fatalln("Error on panamax buffer exec:", err)
		return
	}

	log.Printf("from NRPE Response: %#v", resp)
}

// Extract command and arguments from the packet buffer and compose a call
// towards Panamax.
func (ns *NRPESrvc) panamaxExec(cmd string, args []string) (*Response, error) {
	var err error
	var resp *Response

	log.Printf("About to Panamax Execute: cmd: '%s', args: '%s'", cmd, args)

	if resp, err = ns.p.Execute(cmd, args); err != nil {
		log.Fatalln("Error on Panamax.Execute:", err)
		return nil, err
	}

	return resp, nil
}

func (ns *NRPESrvc) Stop() {
	var err error
	if err = ns.listener.Close(); err != nil {
		log.Fatalln("Error on closing listener:", err)
	}
}

/* EOF */
