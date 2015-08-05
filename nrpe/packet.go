package nrpe

/*
// ---------------------------------------------------------------------------
// Skeleton data structure used by NRPE query/response packages, code is
// inspired on their original source code (version 2.15), which can be found
// on Source-Forge:
//     http://downloads.sourceforge.net/project/nagios/nrpe-2.x/nrpe-2.15
//

#include <stdlib.h>
#include <arpa/inet.h>
#include <sys/socket.h>

// maximum size of a query/response buffer
#define MAX_PACKETBUFFER_LENGTH	1024

typedef struct packet_struct {
	int16_t   packet_version;
	int16_t   packet_type;
	u_int32_t crc32_value;
	int16_t   result_code;
	char      buffer[MAX_PACKETBUFFER_LENGTH];
} packet;

// EOF
*/
import "C"

import (
	"encoding/binary"
	"log"
	// "net"
	"bytes"
	"fmt"
	"unsafe"
    "errors"
)

const (
	NRPE_PACKET_VERSION_3 = 3
	NRPE_PACKET_VERSION_2 = 2
	NRPE_PACKET_VERSION_1 = 1
	NRPE_PACKET_QUERY     = 1
	NRPE_PACKET_RESPONSE  = 2
	NRPE_PACKET_SIZE      = 1036
	NRPE_HELLO_COMMAND    = "_NRPE_CHECK"

	MAX_PACKETBUFFER_LENGTH = 1024
	MAX_COMMAND_ARGUMENTS   = 16

	DEFAULT_SOCKET_TIMEOUT     = 10
	DEFAULT_CONNECTION_TIMEOUT = 300

	STATE_UNKNOWN  = 3
	STATE_CRITICAL = 2
	STATE_WARNING  = 1
	STATE_OK       = 0
)

// Translation of nrpe original data structure into Go standards
type Packet struct {
	PacketVersion int16
	PacketType    int16
	CRC32Value    uint32
	Buffer        string
}

// Extracts a NRPE package from a net.Conn typed argument, it's tranlated into
// oringinal transport format (*C.char) and then binary converted to have a
// local CGO struct.
func Unassemble(cbytes []byte, size int) (pkt *Packet, err error) {
    // checking for the packet size read from socket
    if size != NRPE_PACKET_SIZE {
        __fatal_message := fmt.Sprintf(
            "Wrong packet size: %d/%d\n", size, NRPE_PACKET_SIZE)
        log.Fatalf(__fatal_message)
        return nil, errors.New(__fatal_message)
    }

	// allocating a array of C.char
	c_char := make([]C.char, NRPE_PACKET_SIZE)
	buffer := bytes.NewBuffer(cbytes)

	// reading the binary provided by connection
	err = binary.Read(buffer, binary.BigEndian, c_char)
	if err != nil {
		log.Fatalln("Failed on binary.Read:", err.Error())
		return nil, err
	}

	// casting extracted c.char array into a packet struct
	c_packet := (*C.packet)(unsafe.Pointer(&c_char[0]))
	// special treatment for "buffer" packet entry, based on C.char array
	c_packet_buffer := (*C.char)(unsafe.Pointer(&c_packet.buffer))

	// creating a new go struct to represent NRPE packet and casting to
	// convert into Go types, and "hlons(3)" for network byte order
	go_packet := &Packet{
		PacketVersion: (int16)(C.htons((C.uint16_t)(c_packet.packet_version))),
		PacketType:    (int16)(C.htons((C.uint16_t)(c_packet.packet_type))),
		CRC32Value:    (uint32)(C.htonl((C.uint32_t)(c_packet.crc32_value))),
		Buffer:        (string)(C.GoString(c_packet_buffer)),
	}

	return go_packet, nil
}

/* EOF */
