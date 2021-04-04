package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/hrabalvojta/grpc-test/gateway"
	"github.com/hrabalvojta/grpc-test/insecure"
	pbExample "github.com/hrabalvojta/grpc-test/proto"
	"github.com/hrabalvojta/grpc-test/server"
)

// Command line defaults
const (
	DefaultgRPCAddr           = "127.0.0.1:10000"
	DefaultSwaggerAddr        = "127.0.0.1:11000"
	DefaultPostgreSQLAddr     = "127.0.0.1:5432"
	DefaultPostgreSQLDB       = "dvdrental"
	DefaultPostgreSQLUser     = "postgres"
	DefaultPostgreSQLPassword = "secret"
)

// Command line parameters
var gRPCAddr string
var swaggerAddr string
var pgAddr string
var pgDB string
var pgUser string
var pgPassword string

func init() {
	flag.StringVar(&gRPCAddr, "grpc-addr", DefaultgRPCAddr, "Set the gRPC bind address")
	flag.StringVar(&swaggerAddr, "swagger-addr", DefaultSwaggerAddr, "Set the gRPC bind address")
	flag.StringVar(&pgAddr, "pg-addr", DefaultPostgreSQLAddr, "Set PostgreSQL address")
	flag.StringVar(&pgDB, "pg-db", DefaultPostgreSQLDB, "Set PostgreSQL database")
	flag.StringVar(&pgUser, "pg-user", DefaultPostgreSQLUser, "Set PostgreSQL user")
	flag.StringVar(&pgPassword, "pg-password", DefaultPostgreSQLPassword, "Set PostgreSQL password")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	lis, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer(
		// TODO: Replace with your own certificate!
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
	)
	pbExample.RegisterUserServiceServer(s, server.New())

	// Serve gRPC Server
	log.Info("[main] Serving gRPC on https://", gRPCAddr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	err = gateway.Run("dns:///"+gRPCAddr, swaggerAddr)
	log.Fatalln(err)
}
