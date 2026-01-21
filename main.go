package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)


func main() {
	serverIP := findServer()
	fmt.Println("connecting to: ", serverIP, "...")

	conn, err := net.Dial("tcp",serverIP+ ":8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	sendMessage(conn)
	go handleConnection(conn)

	conn.Close()
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
			os.Exit(1)
		}
		fmt.Print(message)
	}
}

func findServer() string {
	fmt.Println("Looking for server...")
	addr, _ := net.ResolveUDPAddr("udp", ":9999")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening for beacon:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	_, remoteAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error reading from UDP: ", err)
	}

	fmt.Println("Found server at:", remoteAddr.IP.String())
    return remoteAddr.IP.String()
}


