package main

/*
#include <stdlib.h>
#include "common.h"
#include "nrpe.h"

packet *cgo_nrpe_packet(char *recv_packet) {
    return (packet *) recv_packet;
}
*/
import "C"

import (
	"fmt"
	"github.com/otaviof/godutch/config"
	// "github.com/otaviof/godutch/nrpe"
	"bytes"
	"encoding/binary"
	// "io"
	"net"
	"os"
	"unsafe"
	// "reflect"
)

type Obj C.packet

func to_c_struct_packet(cbytes []byte) (result *_Ctype_packet) {
	buf := bytes.NewBuffer(cbytes)
	// var c_packet *_Ctype_char = buf
	fmt.Println("Debug -> buf: ", buf)
	// return (*_Ctype_packet)(*C.char)(&buf)
	//
	// var ptr C.char
	// if err := binary.Read(buf, binary.BigEndian, &ptr); err == nil {
	//    fmt.Println("Debug -> ptr: ", ptr)
	//    return (*_Ctype_packet)(unsafe.Pointer(&ptr))
	// }
	return nil
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	// request := make([]byte, 1036)
	c_char := make([]C.char, 1036)

	// Read the incoming connection into the buffer.
	// _, err := io.ReadFull(conn, c_char)
	err := binary.Read(conn, binary.BigEndian, c_char)
	if err != nil {
		panic("Error reading: " + err.Error())
	}
	fmt.Println("Debug -> c_char #", c_char, "#")

	c_struct := (*C.packet)(unsafe.Pointer(&c_char[0]))
	fmt.Println("Debug -> c_struct (packet_version) #", (int)(c_struct.packet_version), "#")
	fmt.Printf("Debug -> struct #%+v#\n", c_struct)

	c_buffer := (*C.char)(unsafe.Pointer(&c_struct.buffer))
	fmt.Println("Debug -> c_buffer #", c_buffer, "#")
	fmt.Println("Debug -> c_buffer #", C.GoString(c_buffer), "#")

	// go_buffer := (*C.char)(unsafe.Pointer(&buffer))
	// fmt.Println("Debug -> go_buffer #", go_buffer, "#")
	// fmt.Println("Debug -> go_buffer #", C.GoString(go_buffer), "#")

	// pkt := (*Obj)(to_c_struct_packet(request))
	// fmt.Println("Debug -> pkt: ", pkt.buffer)
	// fmt.Println("Debug -> pkt.buffer: ", (* C.struct_packet)(pkt).packet_version)

	// trying to identify which type is the imported buffer
	// value := reflect.ValueOf(pkt)
	// typ := reflect.TypeOf(pkt)

	// fmt.Println("Debug -> value: ", value)
	// fmt.Println("Debug -> typ: ", typ)

	/*
		var c_packet C.struct_packet
		err = binary.Read(buffer, binary.BigEndian, &c_packet)
		if err != nil {
			panic("Error on binary read: " + err.Error())
		}
		fmt.Println("Debug -> C_packet: ", c_packet)

		type GoPacket C.struct_packet
		// var pkt GoPacket
		pkt := (*C.struct_packet)(unsafe.Pointer(&c_packet))
		fmt.Println(pkt)
	*/

	// fmt.Println("C_packet.packet_version: ", (int)(c_packet.packet_version))

	// var pkt *GoPacket
	// pkt := (*GoPacket)(*C.typedef_Packet)(&c_packet)
	// fmt.Println(pkt)

	// Send a response back to person contacting us.
	// conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}

func main() {
	cfg := config.LoadConfig("./etc/godutch/godutch.ini")
	listen_on := fmt.Sprintf(
		"%s:%d", cfg.GoDutch.ListenAddress, cfg.GoDutch.ListenPort)

	l, err := net.Listen("tcp", listen_on)
	if err != nil {
		fmt.Println("Error listening on:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on: " + listen_on)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

/* EOF */
