package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	okStatus       = []byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n")
	notFoundStatus = []byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n")
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil || !strings.HasPrefix(string(buf), "GET / HTTP/1.1") {
		conn.Write(notFoundStatus)
		return
	}

	conn.Write(okStatus)
}

func main() {
	const addr = "0.0.0.0:4221"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConn(conn)
	}
}
