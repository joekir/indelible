package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/joekir/indelible/servicepb"
	"google.golang.org/grpc"
)

const (
	logFilePath = "/var/log/immutable.log"
	socketPath  = "/tmp/indelible.sock"
)

func unixConnect(socketFilePath string, t time.Duration) (net.Conn, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", socketFilePath)
	if err != nil {
		log.Fatalf("Unable to resolve unix-socket address: %v", err)
	}

	conn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		log.Fatalf("Cannot dial unix-socket address: %v", err)
	}
	return conn, err
}

func main() {
	fmt.Printf("Creating log file at %s\n", logFilePath)
	file, err := os.Create(logFilePath)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	conn, err := grpc.Dial(socketPath, grpc.WithInsecure(), grpc.WithDialer(unixConnect))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewLogCreatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Printf("Requesting log file (%s) be marked append-only...\n", logFilePath)
	_, err = c.CreateLog(ctx, &pb.LogFileRequest{Path: logFilePath})
	if err != nil {
		log.Fatalf("Unable to create logfile: %v", err)
	}
	fmt.Println("success")
}
