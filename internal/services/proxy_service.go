package services

import (
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

	_, err = io.Copy(ps.responseWriter, resp.Body)
	if err != nil {
		// log error
	}

	return nil
}

func (ps *ProxyService) handleHttps() error {
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

	go func() {
		defer destConn.Close()
		defer clientConn.Close()
		io.Copy(destConn, clientConn)
	}()
	go func() {
		defer destConn.Close()
		defer clientConn.Close()
		io.Copy(clientConn, destConn)
	}()

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
