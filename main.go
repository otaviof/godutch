package main

import (
	"fmt"
	"github.com/otaviof/godutch/config"
	"github.com/otaviof/godutch/nrpe"
	"net"
	"os"
)

func main() {
	cfg := config.LoadConfig("./__etc/godutch/godutch.ini")
	listen_on := fmt.Sprintf(
		"%s:%d",
		cfg.GoDutch.ListenAddress,
		cfg.GoDutch.ListenPort)

	l, err := net.Listen("tcp", listen_on)
	if err != nil {
		fmt.Println("Error listening on:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on: " + listen_on)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	nrpe_packet, _ := nrpe.Unassemble(conn)
	fmt.Printf("Debug -> nrpe_packet #%+v#\n", nrpe_packet)
	conn.Close()
}

/* EOF */
