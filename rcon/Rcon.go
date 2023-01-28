package rcon

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const SERVERDATA_AUTH = 3
const SERVERDATA_AUTH_RESPONSE = 2
const SERVERDATA_EXECCOMMAND = 2
const SERVERDATA_RESPONSE_VALUE = 0

type Rcon struct {
	Conn     net.Conn
	Host     string
	Port     string
	Password string
}

func (rcon *Rcon) Connect() (err error) {
	rcon.Conn, err = net.DialTimeout("tcp", rcon.Host+":"+rcon.Port, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	// how to deal with unsuccessful init of Rcon?
	// defer rcon.Conn.Close() ??

	// change seed
	rand.Seed(time.Now().UnixNano())

	id := rand.Int31()
	fmt.Println("Authenticating with id:", id)

	p_auth := make_packet(id, SERVERDATA_AUTH, rcon.Password)

	fmt.Println("Auth packet:", p_auth)

	_, err = rcon.Conn.Write(p_auth)
	if err != nil {
		fmt.Println(err)
		return
	}

	// read auth response
	auth_response_buffer := make([]byte, 4096)
	len, _ := rcon.Conn.Read(auth_response_buffer)

	fmt.Println("Auth response:", auth_response_buffer[:len])

	// unpack auth response
	auth_id, auth_response_type, _ := unpack_packet(auth_response_buffer)
	if auth_response_type != SERVERDATA_AUTH_RESPONSE {
		fmt.Println("Unknown auth response type:", auth_response_type)
		return
	}

	if auth_id == -1 {
		fmt.Println("Auth failed")
		return
	}

	fmt.Println("Connected to", rcon.Host+":"+rcon.Port)

	return
}

func make_packet(id int32, _type int, data string) []byte {
	packet := make([]byte, 0)

	// any better solutions?

	// id
	p_id := make([]byte, 4)
	binary.LittleEndian.PutUint32(p_id, uint32(id))
	packet = append(packet, p_id...)

	// type
	p_type := make([]byte, 4)
	binary.LittleEndian.PutUint32(p_type, uint32(_type))
	packet = append(packet, p_type...)

	// data
	p_data := []byte(data)
	p_data = append(p_data, 0x00)
	packet = append(packet, p_data...)

	// length
	p_length := make([]byte, 4)
	binary.LittleEndian.PutUint32(p_length, uint32(len(packet)))
	packet = append(p_length, packet...)

	fmt.Printf("Packet id: %v", p_id)

	return packet
}

func unpack_packet(packet []byte) (id int, _type int, data string) {

	id = int(binary.LittleEndian.Uint32(packet[4:8]))
	_type = int(binary.LittleEndian.Uint32(packet[8:12]))
	data = string(packet[8:])

	return
}
