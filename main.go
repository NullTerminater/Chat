package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"time"
)


func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
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
	count := 0
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("eror reading:", err)
			count++
			time.Sleep(time.Second * 3)
			if count == 3 {
				err = conn.Close()
				os.Exit(1)
			}
		}
		fmt.Print(message)
	}
}

