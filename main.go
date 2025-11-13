package main

import (
	"fmt"
	_ "fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	listen, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for  connections
	conn, err := listen.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close() // Close connection once finished

	for {

		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		// Ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}

}
