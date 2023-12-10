package middlewares

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type responseData struct {
	RespSize   int
	StatusCode int
}
type loggerRW struct {
	http.ResponseWriter
	Response *responseData
}

func (lrw *loggerRW) Write(buf []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(buf)
	lrw.Response.RespSize = size
	return size, err
}

func (lrw *loggerRW) WriteHeader(statusCode int) {
	lrw.Response.StatusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		uri := r.RequestURI
		start := time.Now()

		var (
			data = &responseData{}
			lrw  = &loggerRW{ResponseWriter: w, Response: data}
		)
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		log.Info().
			Str("method", method).
			Str("url", uri).
			Msgf("Request:")
		//Msgf("Request: %s %s time: %fs",
		//	method, uri, duration.Seconds())

		log.Info().
			Int("status", lrw.Response.StatusCode).
			Int("resp-size", lrw.Response.RespSize).
			Dur("duration", duration).
			Msg("Response:")
	}
	return http.HandlerFunc(logFunc)
}
