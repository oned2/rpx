package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type RPXProxy struct {
	proxy *httputil.ReverseProxy
}

// NewProxy takes target host and creates a reverse proxy
func NewRPXProxy(targetHost string) (*RPXProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	rpx := new(RPXProxy)
	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		rpx.modifyRequest(req)
	}

	proxy.ModifyResponse = rpx.modifyResponse()
	proxy.ErrorHandler = rpx.errorHandler()
	rpx.proxy = proxy
	return rpx, nil
}

func (rpx *RPXProxy) modifyRequest(req *http.Request) {
	// hook to modify request
	// this is where we can add custom login logic
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}

func (rpx *RPXProxy) errorHandler() func(http.ResponseWriter, *http.Request, error) {
	// hook to handle error and log it
	return func(w http.ResponseWriter, req *http.Request, err error) {
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadGateway)
		}
	}
}

func (rpx *RPXProxy) modifyResponse() func(*http.Response) error {
	// hook to modify response
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
		return nil
	}
}

func (rpx *RPXProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// let base proxy handle the request
	rpx.proxy.ServeHTTP(w, r)
}

func ProxyServe(sourceHostPort, targetHostPort string) {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewRPXProxy(targetHostPort)
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	log.Fatal(http.ListenAndServe(sourceHostPort, proxy))
}
