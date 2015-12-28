package godutch

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type NRPEService struct {
	listener net.Listener
	cfg      *NRPEConfig
	p        *Panamax
}

func NewNRPEService(cfg *NRPEConfig, p *Panamax) (*NRPEService, error) {
	var err error
	var ns *NRPEService
	ns = &NRPEService{
		cfg: cfg,
		p:   p,
	}
	return ns, err
}

func (ns *NRPEService) Serve() {
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

func (ns *NRPEService) handleConnection(conn net.Conn) {
	var err error
	var n int
	var buf []byte = make([]byte, NRPE_PACKET_SIZE)
	var pkt *NRPEPacket
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

	// using buffer to exectract command and it's argument
	if resp, err = ns.panamaxExecute(pkt.Buffer); err != nil {
		log.Fatalln("Error on panamax buffer exec:", err)
		return
	}

	log.Printf("from NRPE Response: %#v", resp)
}

func extractCommandFromBuffer(pktBuffer string) (string, []string, error) {
	var command string
	var args []string = []string{}

	// splitting informed buffer based on exclamation marks, defualt for NRPE
	buffer = strings.Split(pktBuffer, "!")
	log.Println("Extracted from NPRE Packet buffer:", buffer[:])

	// checking how many items we have, at least one to compose a command
	switch len(buffer) {
	case 0:
		err = errors.New("Can't extract command from buffer:" + pktBuffer)
		return nil, err
	case 1:
		command = fmt.Sprintf("%s", buffer[0])
	default:
		command = fmt.Sprintf("%s", buffer[0])
		args = buffer[1:]
	}
}

// Extract command and arguments from the packet buffer and compose a call
// towards Panamax.
func (ns *NRPEService) panamaxExecute(command string) (*Response, error) {
	var err error
	var buffer []string
	var resp *Response

	log.Println("Panamax Execute with command: '", command, "', args:", args)

	// redirecting request towards Panamax
	if resp, err = ns.p.Execute(command, args); err != nil {
		log.Fatalln("Error on Panamax.Execute:", err)
		return nil, err
	}

	return resp, nil
}

func (ns *NRPEService) Stop() {
	return
}

/* EOF */
