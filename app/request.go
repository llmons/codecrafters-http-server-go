package main

import (
	"strconv"
	"strings"
)

type Request struct {
	method     string
	target     string
	version    string
	reqHeaders ReqHeaders
	body       string
}

type ReqHeaders struct {
	host          string
	userAgent     string
	accept        string
	contentType   string
	contentLength int
}

func parseRequest(req string) Request {
	request := Request{}
	lines := strings.Split(req, "\r\n")

	// request line
	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")
	request.method = parts[0]
	request.target = parts[1]
	request.version = parts[2]

	// headers
	reqHeaders := lines[1 : len(lines)-1]
	mp := map[string]string{}
	for _, header := range reqHeaders {
		parts := strings.Split(header, ": ")
		if len(parts) == 2 {
			mp[parts[0]] = parts[1]
		}
	}

	if val, ok := mp["Host"]; ok {
		request.reqHeaders.host = val
	}
	if val, ok := mp["User-Agent"]; ok {
		request.reqHeaders.userAgent = val
	}
	if val, ok := mp["Accept"]; ok {
		request.reqHeaders.accept = val
	}
	if val, ok := mp["Content-Type"]; ok {
		request.reqHeaders.contentType = val
	}
	if val, ok := mp["Content-Length"]; ok {
		contentLength, err := strconv.Atoi(val)
		if err == nil {
			request.reqHeaders.contentLength = contentLength
		}
	}

	// body
	body := lines[len(lines)-1:]
	request.body = strings.Trim(body[0], "\x00")

	return request
}
