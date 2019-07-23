package main

import "github.com/spf13/pflag"

const (
	defaultServerPort         = "9090"
	defaultBookServiceAddress = "svc-books:9091"
)

var (
	flagServerPort         = pflag.String("server.port", defaultServerPort, "default server port")
	flagBookServiceAddress = pflag.String("books.service", defaultBookServiceAddress, "default book service address")
)
