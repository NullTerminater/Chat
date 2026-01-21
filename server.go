package main
import (
	"fmt"
	"net"
	"os"
	"bufio"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("error listening", err)
		os.Exit(1)
	}
		conn, err := ln.Accept()
		go handleConnection(conn)
		go broadcast()
	for {
		sendMessage(conn)
		if err != nil {
			fmt.Println("erorr accepting: ", err)
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


func broadcast() {
	pc, err := net.Dial("udp", "255.255.255.255:9999")
	if err != nil {
		fmt.Println("Error broadcasting", err)
		return
	}
	defer pc.Close()
	for {
		_, err = pc.Write([]byte("Here"))
		if err != nil {
			fmt.Println("Broadcast error:", err)

		}
		time.Sleep(time.Second * 1)
	}
}

