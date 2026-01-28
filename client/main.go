package main

import (
	"fmt"
	"net"
	"os"

	"terminalcomm/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	serverPort    = ":8080"
	discoveryPort = ":9999"
)

func main() {
	// 1. Discovery Phase (Blocking, before TUI starts)
	serverIP := findServer()

	serverAddr, err := net.ResolveUDPAddr("udp", serverIP+serverPort)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		os.Exit(1)
	}

	// Connect (DialUDP)
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 2. Initialize Shared UI Model
	// Define the callback function for sending messages
	sendCallback := func(msg string) error {
		_, err := conn.Write([]byte(msg))
		return err
	}

	model := ui.InitialModel("UDP Chat Client", sendCallback)

	// 3. Start TUI
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Start network listener in a goroutine
	go func() {
		buf := make([]byte, 2048)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				// If read fails (e.g. conn closed), maybe quit program or log
				return
			}
			text := string(buf[:n])
			// Send message to Bubble Tea program
			p.Send(ui.UdpMsg(fmt.Sprintf("[Server]: %s", text)))
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func findServer() string {
	fmt.Println("Listening for server beacon on UDP :9999...")
	addr, err := net.ResolveUDPAddr("udp", discoveryPort)
	if err != nil {
		fmt.Println("Address resolve error:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Listen error (is another client running?):", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	_, remoteAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Read error:", err)
		os.Exit(1)
	}

	return remoteAddr.IP.String()
}
