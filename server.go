package main
import (
	"fmt"
	"net"
	"os"
	"bufio"
	// "time"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("error listening", err)
		os.Exit(1)
	}
	for{
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("erorr accepting: ", err)
		}
		go handleConnection(conn)
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
