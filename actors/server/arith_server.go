package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	arith "example.com/m/arith"
)

func arith_server(ctx context.Context, rwc io.ReadWriteCloser) error {
	// Instantiate a local ArithServer.
	server := arith.ArithServer{}

	// Derive a client capability that points to the server.  Note the
	// return type of arith.ServerToClient.  It is of type arith.Arith,
	// which is the client capability.  This capability is bound to the
	// server instance above; calling client methods will result in RPC
	// against the corresponding server method.
	//
	// The client can be shared over the network.
	client := arith.Arith_ServerToClient(server)

	// Expose the client over the network.  The 'rwc' parameter can be any
	// io.ReadWriteCloser.  In practice, it is almost always a net.Conn.
	//
	// Note the BootstrapClient option.  This tells the RPC connection to
	// immediately make the supplied client -- an arith.Arith, in our case
	// -- to the remote endpoint.  The capability that an rpc.Conn exports
	// by default is called the "bootstrap capability".
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
		// The BootstrapClient is the RPC interface that will be made available
		// to the remote endpoint by default.  In this case, Arith.
		BootstrapClient: capnp.Client(client),
	})
	defer conn.Close()

	// Block until the connection terminates.
	select {
	case <-conn.Done():
		return nil
	case <-ctx.Done():
		return conn.Close()
	}
}

func accepting_req(ctx context.Context, l net.Listener) {
	c1, err := l.Accept()
	if err != nil {
		panic(err)
	}

	if err := arith_server(ctx, c1); err != nil {
		log.Println(err)
	}
}

func main() {
	const SOCK_ADDR = "./target/example.sock"
	err := os.RemoveAll(SOCK_ADDR)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("unix", SOCK_ADDR)
	if err != nil {
		panic(err)
	}

	defer l.Close()
	ctx := context.Background()

	for {
		accepting_req(ctx, l)
	}

}
