package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverIP := findServer()
	fmt.Printf("Connecting to server at %s on UDP :8080...\n", serverIP)

	serverAddr, err := net.ResolveUDPAddr("udp", serverIP+":8080")
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("--- UDP Chat Client ---")
	fmt.Println("Type your message and press Enter. Receiving is active.")

	// Start a goroutine to receive messages from the server
	go func() {
		buf := make([]byte, 2048)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			// Use \r and clear line if possible, or just print
			fmt.Printf("\r[Server]: %s> ", string(buf[:n]))
		}
	}()

	// Main loop to send messages to the server
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			continue
		}
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}

func findServer() string {
	fmt.Println("Listening for server beacon on UDP :9999...")
	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		fmt.Println("Error resolving discovery address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening for beacon:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error reading from UDP:", err)
		os.Exit(1)
	}

	fmt.Printf("Found server at %s. (Beacon: %s)\n", remoteAddr.IP.String(), string(buf[:n]))
	return remoteAddr.IP.String()
}
