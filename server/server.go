package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("--- UDP Chat Server ---")
	fmt.Println("Listening on UDP :8080")

	go broadcast()

	var clientAddr *net.UDPAddr
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			if clientAddr != nil {
				_, err := conn.WriteToUDP([]byte(text), clientAddr)
				if err != nil {
					fmt.Println("Error sending to client:", err)
				}
			} else {
				fmt.Println("No client discovered yet. Waiting for a message from client...")
			}
			fmt.Print("> ")
		}
	}()

	buf := make([]byte, 2048)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		clientAddr = remoteAddr

		fmt.Printf("\r[Client %s]: %s> ", remoteAddr.String(), string(buf[:n]))
	}
}

func broadcast() {
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	if err != nil {
		fmt.Println("Broadcast resolve error:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Broadcast dial error:", err)
		return
	}
	defer conn.Close()

	for {
		_, err = conn.Write([]byte("SERVER_ALIVE"))
		if err != nil {
		}
		time.Sleep(2 * time.Second)
	}
}
