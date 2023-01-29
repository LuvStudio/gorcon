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
	conn     net.Conn
	Host     string
	Port     string
	Password string
}

func (rcon *Rcon) Connect() (err error) {
	rcon.conn, err = net.DialTimeout("tcp", rcon.Host+":"+rcon.Port, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	// how to deal with unsuccessful init of Rcon?
	// defer rcon.conn.Close() ??

	// change seed
	rand.Seed(time.Now().UnixNano())

	id := rand.Int31()

	p_auth := make_packet(id, SERVERDATA_AUTH, rcon.Password)

	_, err = rcon.conn.Write(p_auth)
	if err != nil {
		fmt.Println(err)
		return
	}

	// read auth response
	auth_response_buffer := make([]byte, 4096)
	len, err := rcon.conn.Read(auth_response_buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// unpack auth response
	auth_id, auth_response_type, _ := unpack_packet(auth_response_buffer[:len])
	if auth_response_type != SERVERDATA_AUTH_RESPONSE {
		return fmt.Errorf("unknown auth response type: %v", auth_response_type)
	}

	if auth_id == -1 {
		return fmt.Errorf("auth failed")
	}

	fmt.Println("Connected to", rcon.Host+":"+rcon.Port)

	return
}

func (rcon *Rcon) Close() (err error) {
	print("Closing connection to ", rcon.Host+":"+rcon.Port)
	err = rcon.conn.Close()
	return
}

func (rcon *Rcon) Command(data string) (response string, err error) {

	rand.Seed(time.Now().UnixNano())
	id := rand.Int31()

	p_command := make_packet(id, SERVERDATA_EXECCOMMAND, data)

	_, err = rcon.conn.Write(p_command)
	if err != nil {
		return "", err
	}

	// read command response
	command_response_buffer := make([]byte, 4096)
	len, err := rcon.conn.Read(command_response_buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// unpack command response
	command_id, command_response_type, command_response_data := unpack_packet(command_response_buffer[:len])
	if command_response_type != SERVERDATA_RESPONSE_VALUE {
		return "", fmt.Errorf("unknown command response type: %v", command_response_type)
	}

	if command_id != id {
		return "", fmt.Errorf("unknown command id: %v", command_id)
	}

	return command_response_data, nil
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

	return packet
}

func unpack_packet(packet []byte) (id int32, _type int, data string) {

	id = int32(binary.LittleEndian.Uint32(packet[4:8]))
	_type = int(binary.LittleEndian.Uint32(packet[8:12]))
	data = string(packet[8:])

	return
}
