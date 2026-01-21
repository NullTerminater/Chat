package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)


func main() {
	conn, err := net.Dial("tcp", ":8080")
	for {
		sendMessage(conn)
	}
	err = conn.Close()
	if err != nil {
		fmt.Println("Error closing", err)
	}
}

func sendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadString('\n')
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending the message", err)
			break
		}
	}
}
