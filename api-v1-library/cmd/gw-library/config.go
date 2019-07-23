package main

import "github.com/spf13/pflag"

const (
	defaultGatewayPort    = "8081"
	defaultServerPort     = "9090"
	defaultServerEndpoint = "svc-library:9090"
)

var (
	flagGatewayPort    = pflag.String("gateway.port", defaultGatewayPort, "default gateway port")
	flagServerPort     = pflag.String("server.port", defaultServerPort, "default server port")
	flagServerEndpoint = pflag.String("server.endpoint", defaultServerEndpoint, "default server endpoint")
)
