package services

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"krizanauskas.github.com/mvp-proxy/internal/errors"
)

type ProxyService struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

var httpsCopyDone = make(chan struct{}, 2)

var httpsCopyFunc = func(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Printf("Error copying data: %v", err)
	}
	httpsCopyDone <- struct{}{}
}

var httpCopyDone = make(chan error, 1)

var httpCopyFunc = func(dst http.ResponseWriter, resp io.ReadCloser) {
	_, copyErr := io.Copy(dst, resp)
	httpCopyDone <- copyErr
}

func NewProxyService(w http.ResponseWriter, r *http.Request) *ProxyService {
	return &ProxyService{
		w,
		r,
	}
}

func (ps *ProxyService) ProxyRequest() error {
	if ps.request.Method == http.MethodConnect {
		return ps.handleHttps()
	}

	return ps.handleHttp()
}

func (ps *ProxyService) handleHttp() error {
	ctx := ps.request.Context()

	proxyReq, err := http.NewRequest(ps.request.Method, ps.request.URL.String(), ps.request.Body)
	if err != nil {
		return &errors.ServiceError{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)}
	}
	proxyReq.Header = ps.request.Header.Clone()

	proxyReq = proxyReq.WithContext(ps.request.Context())

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return &errors.ServiceError{Code: http.StatusBadGateway, Message: http.StatusText(http.StatusBadGateway)}
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		if isHopByHopHeader(key) {
			continue
		}

		for _, value := range values {
			ps.responseWriter.Header().Add(key, value)
		}
	}

	ps.responseWriter.WriteHeader(resp.StatusCode)

	go httpCopyFunc(ps.responseWriter, resp.Body)

	select {
	case <-ctx.Done():
		// Context was canceled or timed out
		fmt.Printf("Context canceled or timed out: %v", ctx.Err())
		return ctx.Err()
	case copyErr := <-httpCopyDone:
		if copyErr != nil {
			fmt.Printf("Error copying response body: %v", copyErr)
			return copyErr
		}
	}

	return nil
}

func (ps *ProxyService) handleHttps() error {
	ctx := ps.request.Context()

	destConn, err := net.Dial("tcp", ps.request.Host)
	if err != nil {
		return &errors.ServiceError{Code: http.StatusServiceUnavailable, Message: http.StatusText(http.StatusServiceUnavailable)}
	}

	ps.responseWriter.WriteHeader(http.StatusOK)

	hijacker, ok := ps.responseWriter.(http.Hijacker)
	if !ok {
		return &errors.ServiceError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
	}

	// hijack client TCP connection
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		return &errors.ServiceError{Code: http.StatusServiceUnavailable, Message: http.StatusText(http.StatusServiceUnavailable)}
	}

	go httpsCopyFunc(destConn, clientConn)
	go httpsCopyFunc(clientConn, destConn)

	// Listen for context cancellation or copying completion
	select {
	case <-ctx.Done():
		fmt.Printf("Context canceled or timed out: %v", ctx.Err())
		// Close both connections to terminate the tunnel
		clientConn.Close()
		destConn.Close()
		return ctx.Err()
	case <-httpsCopyDone:
	}

	<-httpsCopyDone

	fmt.Println("finished")

	return nil
}

func isHopByHopHeader(header string) bool {
	header = http.CanonicalHeaderKey(header)
	for _, h := range hopHeaders {
		if h == header {
			return true
		}
	}
	return false
}
