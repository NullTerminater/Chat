package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)


func main() {
	conn, err := net.Dial("tcp", ":8080")
	go sendMessage(conn)
	for {
		handleConnection(conn)
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

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("eror reading:", err)
			err = conn.Close()
			break
		}
		fmt.Print(message)
	}
}

