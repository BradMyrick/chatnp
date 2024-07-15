package main

import (
	"context"
	"log"
	"net"

	"capnproto.org/go/capnp/v3/rpc"
	"github.com/BradMyrick/chatnp/server"
)

func main() {
    client, srv := server.NewServer()
    defer srv.Shutdown()

    ctx := context.Background()

    log.Println("Starting server...")
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Error accepting connection:", err)
            continue
        }

        rpcConn := rpc.NewConn(rpc.NewStreamTransport(conn), nil)
        go func() {
            select {
            case <-rpcConn.Done():
                client.Release()
            case <-ctx.Done():
                _ = rpcConn.Close()
            }
        }()
    }
}
