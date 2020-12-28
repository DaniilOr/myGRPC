package main

import (
	"context"
	"github.com/DaniilOr/myGRPC/cmd/templates/server/app"
	serverPb "github.com/DaniilOr/myGRPC/pkg/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"
const defaultDSN = "postgres://app:pass@localhost:5432/autopayments"
func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	dsn, ok := os.LookupEnv("APP_HOST")
	if !ok {
		dsn = defaultDSN
	}

	if err := execute(net.JoinHostPort(host, port), dsn); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, dsn string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx := context.Background()
	grpcServer := grpc.NewServer()
	server := app.NewServer(db, ctx)
	serverPb.RegisterPayServiceServer(grpcServer, server)

	return grpcServer.Serve(listener)
}
