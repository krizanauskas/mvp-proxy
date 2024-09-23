package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"krizanauskas.github.com/mvp-proxy/internal/errors"
)

type ProxyService struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	userService    UserServiceI
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

func copyWithBandwidthFunc(ctx context.Context, dst io.Writer, src io.Reader, user string, bandwidthController UserBandwidthControllerI, methodConnect bool) error {
	// 1 MB buffer for copying data in chunks
	const bufferSize = 1000 * 1000

	buffer := make([]byte, bufferSize)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		userLimit := bandwidthController.GetAvailableBandwidth(user)

		if userLimit <= 0 {
			fmt.Printf("No available bandwidth for user: %s\n", user)
			break // Exit the loop if no bandwidth is available
		}

		limitedReader := &io.LimitedReader{
			R: src,
			N: int64(userLimit),
		}

		n, err := limitedReader.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading data: %w", err)
		}

		if n == 0 {
			if err == io.EOF {
				return nil
			}
			break
		}

		if methodConnect {
			written, writeErr := dst.Write(buffer[:n])
			bandwidthController.UpdateBandwidthUsed(user, written)

			if writeErr != nil {
				return fmt.Errorf("error writing data: %w", writeErr)
			}
		} else {
			copied, copyErr := io.CopyN(dst, bytes.NewReader(buffer[:n]), int64(n))
			bandwidthController.UpdateBandwidthUsed(user, int(copied))

			if copyErr != nil {
				return fmt.Errorf("error writing data: %w", copyErr)
			}
		}
	}

	return nil
}

func NewProxyService(w http.ResponseWriter, r *http.Request, userService UserServiceI) *ProxyService {
	return &ProxyService{
		w,
		r,
		userService,
	}
}

func (ps *ProxyService) ProxyRequest() error {
	ps.userService.AddToHistory(AuthUser, getFullURL(ps.request), time.Now())

	if ps.request.Method == http.MethodConnect {
		return ps.handleHttps()
	}

	return ps.handleHttp()
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
		destConn.Close()
		return &errors.ServiceError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
	}

	// hijack client TCP connection
	clientConn, _, err := hijacker.Hijack()
	if err != nil || clientConn == nil {
		destConn.Close()
		return &errors.ServiceError{
			Code:    http.StatusServiceUnavailable,
			Message: "Failed to hijack connection",
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer wg.Done()
		defer clientConn.Close()

		if err := copyWithBandwidthFunc(ctx, destConn, clientConn, AuthUser, ps.userService, true); err != nil {
			fmt.Printf("Error copying from destConn to clientConn: %v\n", err)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		defer destConn.Close()

		if err := copyWithBandwidthFunc(ctx, clientConn, destConn, AuthUser, ps.userService, true); err != nil {
			fmt.Printf("Error copying from clientConn to destConn: %v\n", err)
			cancel()
		}
	}()

	doneChan := make(chan struct{})

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-ctx.Done():
		fmt.Printf("Context canceled or timed out: %v", ctx.Err())

		clientConn.Close()
		destConn.Close()
		return ctx.Err()
	case <-doneChan:
	}

	return nil
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

	var copyErrChan = make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer close(copyErrChan)

		err := copyWithBandwidthFunc(ctx, ps.responseWriter, resp.Body, AuthUser, ps.userService, true)

		if err != nil {
			fmt.Printf("Error copying from destination to client: %v\n", err)
			copyErrChan <- err
		}
	}()

	select {
	case copyErr := <-copyErrChan:
		if copyErr != nil {
			fmt.Printf("Error copying response body: %v", copyErr)
			return copyErr
		}

		return nil
	case <-ctx.Done():
		// Context was canceled or timed out
		fmt.Printf("Context canceled or timed out: %v\n", ctx.Err())
		return ctx.Err()
	}
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

func getFullURL(r *http.Request) string {
	scheme := "http"

	if r.Method == http.MethodConnect {
		scheme = "https"
	}

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}
