package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
		return
	}

	req := parseRequest(string(buf))

	res := Response{}
	res.version = req.version

	if req.target == "/" {
		res.code = 200
		res.status = "OK"
		conn.Write(seqResponse(res))
		return
	}

	if strings.HasPrefix(req.target, "/echo/") {
		res.code = 200
		res.status = "OK"
		res.resHeaders.contentType = "text/plain"
		str := req.target[len("/echo/"):]
		res.resHeaders.contentLength = len(str)
		res.body = str
		conn.Write(seqResponse(res))
		return
	}

	res.code = 404
	res.status = "Not Found"
	conn.Write(seqResponse(res))
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
