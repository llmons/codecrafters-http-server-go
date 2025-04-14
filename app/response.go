package main

import "strconv"

type Response struct {
	version    string
	code       int
	status     string
	resHeaders ResHeaders
	body       string
}

type ResHeaders struct {
	contentType   string
	contentLength int
}

func seqResponse(res Response) []byte {
	ret := []byte{}
	ret = append(ret, []byte(res.version+" "+strconv.Itoa(res.code)+" "+res.status+"\r\n")...)
	ret = append(ret, []byte("Content-Type: "+res.resHeaders.contentType+"\r\n")...)
	ret = append(ret, []byte("Content-Length: "+strconv.Itoa(res.resHeaders.contentLength)+"\r\n")...)
	ret = append(ret, []byte("\r\n")...)
	ret = append(ret, []byte(res.body)...)
	return ret
}
