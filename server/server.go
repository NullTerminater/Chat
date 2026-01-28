package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"terminalcomm/ui"

	tea "github.com/charmbracelet/bubbletea"
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

	go broadcast()

	var clientAddr *net.UDPAddr

	sendCallback := func(msg string) error {
		if clientAddr == nil {
			return fmt.Errorf("no client connected yet")
		}
		_, err := conn.WriteToUDP([]byte(msg), clientAddr)
		return err
	}

	model := ui.InitialModel("UDP Chat Server", sendCallback)
	p := tea.NewProgram(model, tea.WithAltScreen())

	go func() {
		buf := make([]byte, 2048)
		for {
			n, remoteAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				return
			}

			clientAddr = remoteAddr

			text := string(buf[:n])
			p.Send(ui.UdpMsg(fmt.Sprintf("[Client %s]: %s", remoteAddr.IP.String(), text)))
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func broadcast() {
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
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
