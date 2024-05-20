package main

import (
	"context"
	"io"
	"log"
	"net"

	"capnproto.org/go/capnp/v3/rpc"
	arith "example.com/m/arith"
)

func arith_client(ctx context.Context, rwc io.ReadWriteCloser) error {
	// As before, rwc can be any io.ReadWriteCloser, and will typically be
	// a net.Conn.  The rpc.Options can be nil, if you don't want to override
	// the defaults.
	//
	// Here, we expect to receive an arith.Arith from the remote side.  The
	// remote side is not expecting a capability in return, however, so we
	// don't need to define a bootstrap interface.
	//
	// This last point bears emphasis:  capnp RPC is fully bidirectional!  Both
	// sides of a connection MAY export a boostrap interface, and in such cases,
	// the bootstrap interfaces need not be the same!
	//
	// Again, for the avoidance of doubt:  only the remote side is exporting a
	// bootstrap interface in this example.
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), nil)
	defer conn.Close()

	// Now we resolve the bootstrap interface from the remote ArithServer.
	// Thanks to Cap'n Proto's promise pipelining, this function call does
	// NOT block.  We can start making RPC calls with 'a' immediately, and
	// these will transparently resolve when bootstrapping completes.
	//
	// The context can be used to time-out or otherwise abort the bootstrap
	// call.   It is safe to cancel the context after the first method call
	// on 'a' completes.
	a := arith.Arith(conn.Bootstrap(ctx))

	// Okay! Let's make an RPC call!  Remember:  RPC is performed simply by
	// calling a's methods.
	//
	// There are couple of interesting things to note here:
	//  1. We pass a callback function to set parameters on the RPC call.  If the
	//     call takes no arguments, you MAY pass nil.
	//  2. We return a Future type, representing the in-flight RPC call.  As with
	//     the earlier call to Bootstrap, a's methods do not block.  They instead
	//     return a future that eventually resolves with the RPC results. We also
	//     return a release function, which MUST be called when you're done with
	//     the RPC call and its results.
	m, release := a.Multiply(ctx, func(ps arith.Arith_multiply_Params) error {
		ps.SetA(2)
		ps.SetB(32)
		return nil
	})
	defer release()

	d, release := a.Divide(ctx, func(qr arith.Arith_divide_Params) error {
		qr.SetNum(55)
		qr.SetDenom(0)
		return nil
	})
	defer release()

	dres, err := d.Struct()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(dres.Quo(), dres.Rem())
	}

	// You can do other things while the RPC call is in-flight.  Everything
	// is asynchronous. For simplicity, we're going to block until the call
	// completes.
	mres, err := m.Struct()
	if err != nil {
		log.Println(err)
	} else {
		// Lastly, let's print the result.  Recall that 'product' is the name of
		// the return value that we defined in the schema file.
		log.Println(mres.Product()) // prints 84
	}

	return nil
}

func main() {
	const SOCK_ADDR = "./target/example.sock"

	ctx := context.Background()
	c2, err := net.Dial("unix", SOCK_ADDR)
	if err != nil {
		panic(err)
	}

	if err := arith_client(ctx, c2); err != nil {
		log.Println(err)
	}

}
