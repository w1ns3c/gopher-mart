package httpserver

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"sync"
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer(address string, router http.Handler) (srv *HTTPServer, err error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	return &HTTPServer{Server: &http.Server{
		Addr:    addr.String(),
		Handler: router,
	}}, nil

}

func (srv *HTTPServer) Run(ctx context.Context) error {

	log.Info().
		Str("addr", srv.Addr).
		Msg("Server started at:")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		srv.ListenAndServe()
	}()
	<-ctx.Done()
	wg.Done()
	return fmt.Errorf("HTTP server stoped")
}
