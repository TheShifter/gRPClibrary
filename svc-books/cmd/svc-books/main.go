package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"

	pb "github.com/TheShifter/gRPClibrary/svc-books/api/proto"
	_ "github.com/lib/pq"
)

const (
	host     = "db-library"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	if err := RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

type Config struct {
	GRPCPort string
	db       *sql.DB
}

type bookServer struct {
	DB *sql.DB
}

func NewBookServer(db *sql.DB) pb.IBookServer {
	return &bookServer{DB: db}
}

func (s *bookServer) GetBook(ctx context.Context, rq *pb.GetBookRequest) (*pb.Book, error) {
	sqlStatement := `SELECT id, name, description, author, pages FROM books WHERE id=$1;`
	row := s.DB.QueryRow(sqlStatement, rq.Id)
	var book pb.Book
	if err := row.Scan(&book.Id, &book.Name, &book.Description, &book.Author, &book.PageCount); err == sql.ErrNoRows {
		return nil, errors.New("no rows were returned")
	} else if err != nil {
		return nil, err
	}
	return &book, nil
}

func (s *bookServer) ListBooks(ctx context.Context, rq *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	rows, err := s.DB.Query("SELECT id, name, description, author, pages FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []*pb.Book
	for rows.Next() {
		book := pb.Book{}
		err = rows.Scan(&book.Id, &book.Name, &book.Description, &book.Author, &book.PageCount)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &pb.ListBooksResponse{Books: books}, nil
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "9091", "gRPC port to bind")
	flag.Parse()

	//if len(cfg.GRPCPort) == 0 {
	//	cfg.GRPCPort = "9090"
	//return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	//}

	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", info)
	if err != nil {
		panic(err)
	}

	v1API := NewBookServer(db)

	return Start(ctx, v1API, cfg.GRPCPort)
}

func Start(ctx context.Context, v1API pb.IBookServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterIBookServer(server, v1API)

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
