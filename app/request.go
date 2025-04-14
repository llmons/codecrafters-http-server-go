package main

import "strings"

type Request struct {
	method     string
	target     string
	version    string
	reqHeaders ReqHeaders
	body       string
}

type ReqHeaders struct {
	host      string
	userAgent string
	accept    string
}

func parseRequest(req string) Request {
	lines := strings.Split(req, "\r\n")
	requestLine := lines[0]
	tokens := strings.Split(requestLine, " ")
	reqHeaders := lines[1:4]
	body := lines[4:]

	// Request line
	request := Request{}
	request.method = tokens[0]
	request.target = tokens[1]
	request.version = tokens[2]

	// Headers
	if len(reqHeaders[0]) > 0 {
		request.reqHeaders.host = reqHeaders[0][len("Host: "):]
	}
	if len(reqHeaders[1]) > 0 {
		request.reqHeaders.userAgent = reqHeaders[1][len("User-Agent: "):]
	}
	if len(reqHeaders[2]) > 0 {
		request.reqHeaders.accept = reqHeaders[2][len("Accept: "):]
	}

	if len(body) == 0 {
		return request
	}

	// Body
	request.body = body[0]

	return request
}
