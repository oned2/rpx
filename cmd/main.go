package main

import (
	"rpx/proxy"
)

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy.ProxyServe(":8081", "http://localhost:8001")
}
