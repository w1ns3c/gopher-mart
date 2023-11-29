package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// decompress
		if strings.Contains(r.Header.Get("content-encoding"), "gzip") {
			// body compressed
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				// TODO log err
				return
			}
			defer gzipReader.Close()
			r.Body = gzipReader
		}

		// compress
		if strings.Contains(r.Header.Get("accept-encoding"), "gzip") {
			gzipW := &gzipWriter{
				ResponseWriter: w,
				Writer:         gzip.NewWriter(w),
			}
			defer gzipW.Writer.Close()
			gzipW.ResponseWriter.Header().Set("content-encoding", "gzip")
			next.ServeHTTP(gzipW, r)
		}

	})
}
