package main

import (
	"fmt"
	_ "fmt"
	"net"
	"strings"
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

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)

		args := value.array[1:]

		handler, ok := Handlers[command]

		writer := NewWriter(conn)

		if !ok {
			fmt.Println("Invalid command: " + command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)

		writer.Write(result)
	}

}
