package main

import (
	"context"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"net/http"

	gw "github.com/TheShifter/gRPClibrary/api-v1-library/api/proto"
)

//func init() {
//	pflag.Parse()
//	err := viper.BindPFlags(pflag.CommandLine)
//	if err != nil {
//		panic(err)
//	}
//}

func main() {
	pflag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterLibraryHandlerFromEndpoint(ctx, mux, *flagServerEndpoint, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":"+*flagGatewayPort, mux)
}
