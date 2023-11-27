package httpserver

import (
	"net"
	"net/http"
)

type HTTPServer struct {
	addr   *net.TCPAddr
	router *http.ServeMux
}

func NewHTTPServer(address string, router *http.ServeMux) (srv *HTTPServer, err error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	return &HTTPServer{addr: addr, router: router}, nil
}

func (srv *HTTPServer) Run() error {
	return http.ListenAndServe(srv.addr.String(), srv.router)
}
