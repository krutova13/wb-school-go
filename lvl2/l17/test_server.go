package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer listener.Close()

	fmt.Println("Test server started on :8080")
	fmt.Println("Waiting for connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr())

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
			return
		}

		data := string(buffer[:n])
		fmt.Printf("Received from %s: %s", conn.RemoteAddr(), data)

		response := fmt.Sprintf("Echo: %s", strings.TrimSpace(data))
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Printf("Failed to write to client: %v\n", err)
			return
		}
	}
}
