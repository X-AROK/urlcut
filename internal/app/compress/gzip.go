package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var encodingTypes []string = []string{
	"text/html",
	"application/json",
}

type gzipWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func (c *gzipWriter) needsEncoding() bool {
	contentType := c.Header().Get("Content-Type")
	for _, t := range encodingTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}

	return false
}

func (c *gzipWriter) Write(p []byte) (int, error) {
	if !c.needsEncoding() {
		return c.ResponseWriter.Write(p)
	}

	return c.zw.Write(p)
}

func (c *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 && c.needsEncoding() {
		c.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *gzipWriter) Close() error {
	return c.zw.Close()
}

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
