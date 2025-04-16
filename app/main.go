package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleGetReq(req Request) Response {
	res := Response{}
	res.version = req.version
	res.code = 200
	res.status = "OK"

	switch req.target {
	case "/":
	case "/user-agent":
		res.resHeaders.contentType = "text/plain"
		res.resHeaders.contentLength = len(req.reqHeaders.userAgent)
		res.body = req.reqHeaders.userAgent
	default:
		if strings.HasPrefix(req.target, "/echo/") {
			res.resHeaders.contentType = "text/plain"
			str := req.target[len("/echo/"):]
			res.resHeaders.contentLength = len(str)
			res.body = str
		} else if strings.HasPrefix(req.target, "/files/") {
			res.resHeaders.contentType = "application/octet-stream"

			//  read the file from target directory
			dir := os.Args[2]
			path := req.target[len("/files/"):]
			data, err := os.ReadFile(dir + path)
			if err != nil {
				res.code = 404
				res.status = "Not Found"
				return res
			}

			res.resHeaders.contentLength = len(data)
			res.body = string(data)
		} else {
			res.code = 404
			res.status = "Not Found"
		}
	}

	return res
}

func handlePostReq(req Request) Response {
	res := Response{}
	res.version = req.version
	res.code = 201
	res.status = "Created"

	if strings.HasPrefix(req.target, "/files/") {
		// write the file to target directory
		dir := os.Args[2]
		path := req.target[len("/files/"):]
		err := os.WriteFile(dir+path, []byte(req.body), 0644)
		if err != nil {
			res.code = 404
			res.status = "Not Found"
			return res
		}
	} else {
		res.code = 404
		res.status = "Not Found"
	}
	return res
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
		return
	}

	req := parseRequest(string(buf))

	var res Response
	switch req.method {
	case "GET":
		res = handleGetReq(req)
	case "POST":
		res = handlePostReq(req)
	}
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
