package main

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"

	pb "github.com/TheShifter/gRPClibrary/api-v1-library/api/proto"
	pb2 "github.com/TheShifter/gRPClibrary/svc-books/api/proto"
)

func init() {
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

type Config struct {
	GRPCPort string
}

type libraryServer struct {
	BookClient pb2.IBookClient
}

func NewLibraryServer(cl pb2.IBookClient) pb.LibraryServer {
	return &libraryServer{BookClient: cl}
}

func (s *libraryServer) ListBooks(ctx context.Context, rq *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	r := &pb2.ListBooksRequest{
		Name: rq.Name,
	}
	rsp, err := s.BookClient.ListBooks(ctx, r)
	if err != nil {
		return nil, err
	}

	rs := &pb.ListBooksResponse{}
	for _, book := range rsp.Books {
		var retval pb.Book
		retval.Id = book.Id
		retval.Name = book.Name
		retval.Description = book.Description
		retval.Author = book.Author
		retval.PageCount = book.PageCount
		rs.Books = append(rs.Books, &retval)
	}

	return rs, nil
}

func (s *libraryServer) GetBook(ctx context.Context, rq *pb.GetBookRequest) (*pb.Book, error) {
	r := &pb2.GetBookRequest{
		Id: rq.Id,
	}
	rsp, err := s.BookClient.GetBook(ctx, r)
	if err != nil {
		return nil, err
	}

	return &pb.Book{Id: rsp.Id, Name: rsp.Name, Author: rsp.Author, Description: rsp.Description, PageCount: rsp.PageCount}, nil
}

func RunServer() error {
	conn, err := grpc.Dial(*flagBookServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	ctx := context.Background()

	cl := pb2.NewIBookClient(conn)

	v1API := NewLibraryServer(cl)

	return Start(ctx, v1API, *flagServerPort)
}

func Start(ctx context.Context, v1API pb.LibraryServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterLibraryServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
