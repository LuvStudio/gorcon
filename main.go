package main

import (
	"./rcon"
)

const addr = "192.168.0.30"
const port = "25575"

func main() {
	/*
		conn, err := net.Dial("tcp", addr+":"+port)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		fmt.Println("Connected to", addr+":"+port)
		inputReader := bufio.NewReader(os.Stdin)
		for {
			input, _ := inputReader.ReadString('\n')
			conn.Write([]byte(input))
		}
	*/

	Rcon := new(rcon.Rcon)
	Rcon.Host = addr
	Rcon.Port = port
	Rcon.Password = "logitech"

	Rcon.Connect()

	defer Rcon.Conn.Close()

}
