package nrpe

/*
// ---------------------------------------------------------------------------
// Skeleton data structure used by NRPE query/response packages, code is
// inspired on their original source code (version 2.15), which can be found
// on Source-Forge:
//     http://downloads.sourceforge.net/project/nagios/nrpe-2.x/nrpe-2.15
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
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


//
// Returns a array of unsigned long values to dub as CRC32's IEEE table
//
unsigned long *generate_crc32_table(void) {
	unsigned long crc, poly;
	int i, j;
    unsigned long *crc32_table = malloc(sizeof(unsigned long) * 257);

	poly = 0xEDB88320L;

	for (i = 0; i < 256; i++) {
		crc = i;
		for (j = 8; j > 0; j--) {
			if (crc & 1) {
				crc = (crc >> 1) ^ poly;
            } else {
				crc >>= 1;
            }
		}
		crc32_table[i] = crc;
    }


    return crc32_table;
}

//
// Calculates the CRC32 signature of a given C array of Chars, hereby
// represented as it's pointer. The return is the CRC32 unsigned long.
//
unsigned long crc32 (
    char *buffer,
    int buffer_size,
    unsigned long *crc32_table
) {
	register unsigned long crc;
	int this_char;
	int current_index;

	crc = 0xFFFFFFFF;

	for (current_index = 0; current_index < buffer_size; current_index++) {
		this_char = (int)buffer[current_index];
		crc = ((crc >> 8) & 0x00FFFFFF) ^ crc32_table[(crc ^ this_char) & 0xFF];
	}

	return (crc ^ 0xFFFFFFFF);
}

//
// Wrapper method around "crc32", to load a C.packet struct to remove current
// signature and use standard struct to calculate CRC32. The return is a
// unsinged long.
//
unsigned long calc_packet_crc32 (
    packet *receive_packet,
    unsigned long *crc32_table
) {
    unsigned long packet_crc32;
    unsigned long calculated_crc32;
    packet local_packet;

    // copying the received packet to a local variable
    memcpy(&local_packet, receive_packet, sizeof(local_packet));

    // converting back the packet crc32 to u_int32_t
    packet_crc32 = ntohl(local_packet.crc32_value);

    // erasing saved signature so calculcation of struct's CRC32 will have the
    // same keys and values when created
    local_packet.crc32_value = 0L;

	return crc32(
        (char *)&local_packet,
        sizeof(local_packet),
        crc32_table
    );
}
*/
import "C"

import (
	"errors"
	"fmt"
	"log"
	"unsafe"
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

// Translation of NRPE original data structure
type Packet struct {
	PacketVersion int16
	PacketType    int16
	CRC32Value    uint32
	Buffer        string
}

func NewPacket() (pkt *Packet) {
	return nil
}

// Go wrapper around C written method, takes a C.packet as argument and
// returns a uint32 referent to calculated CRC32.
func CalcPacketCRC32(c_packet *C.packet) uint32 {
	c_ieee_table := C.generate_crc32_table()
	crc32_value := (uint32)(C.calc_packet_crc32(c_packet, c_ieee_table))
	log.Println("Calculated CRC32:", crc32_value)
	return crc32_value
}

// Extracts a NRPE package from a byte array typed argument, it's translated
// into original transport format (*C.char) and then binary converted to have
// a local CGO struct returned. Array size (int) is also mandatory as second
// argument.
func Unassemble(cbytes []byte, size int) (pkt *Packet, err error) {
	// checking for the packet size read from socket
	if size != NRPE_PACKET_SIZE {
		__fatal_msg := fmt.Sprintf("Wrong packet size: %d/%d\n",
			size, NRPE_PACKET_SIZE)
		log.Fatalf(__fatal_msg)
		return nil, errors.New(__fatal_msg)
	}

	// extracting original bytes on local C.char pointer
	c_char := (*C.char)(unsafe.Pointer(&cbytes[0]))
	// casting extracted C.char array into a C.packet struct
	c_packet := (*C.packet)(unsafe.Pointer(c_char))
	// validating packet's signature, using bytes content
	c_packet_crc32_value := (uint32)(
		C.ntohl((C.uint32_t)(c_packet.crc32_value)))
	// special treatment for "buffer" packet entry, based on C.char array
	c_packet_buffer := (*C.char)(unsafe.Pointer(&c_packet.buffer))

	// CRC32: validating packat's content via informed CRC32 signature; here
	// using local C methods to integract with raw C.char array
	calculated_crc32 := CalcPacketCRC32(c_packet)

	if c_packet_crc32_value != calculated_crc32 {
		__fatal_msg := fmt.Sprintf(
			"CRC32 mismatch. Calculated: '%ul', packet's: '%ul';\n",
			calculated_crc32, c_packet_crc32_value)
		log.Fatalln(__fatal_msg)
		return nil, errors.New(__fatal_msg)
	}

	// creating a new go struct to represent NRPE packet and casting to
	// convert into Go types, and "hlons(3)" for network byte order
	go_packet := &Packet{
		PacketVersion: (int16)(C.htons((C.uint16_t)(c_packet.packet_version))),
		PacketType:    (int16)(C.htons((C.uint16_t)(c_packet.packet_type))),
		CRC32Value:    (uint32)(c_packet_crc32_value),
		Buffer:        (string)(C.GoString(c_packet_buffer)),
	}

	return go_packet, nil
}

/* EOF */
