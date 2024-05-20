package arith

import (
	context "context"
	"errors"
)

// ArithServer satisfies the Arith_Server interface that was generated
// by the capnp compiler.
type ArithServer struct{}

// Multiply is the concrete implementation of the Multiply method that was
// defined in the schema. Notice that the method signature matches that of
// the Arith_Server interface.
//
// The Arith_multiply struct was generated by the capnp compiler.  You will
// find it in arith.capnp.go
func (ArithServer) Multiply(ctx context.Context, call Arith_multiply) error {
	res, err := call.AllocResults() // allocate the results struct
	if err != nil {
		return err
	}

	// Set the result to be the product of the two arguments, A and B,
	// that we received. These are found in the Arith_multiply struct.
	res.SetProduct(call.Args().A() * call.Args().B())
	return nil
}

// Divide is analogous to Multiply.  All capability server methods follow the
// same pattern.
func (ArithServer) Divide(ctx context.Context, call Arith_divide) error {
	if call.Args().Denom() == 0 {
		return errors.New("divide by zero")
	}

	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	res.SetQuo(call.Args().Num() / call.Args().Denom())
	res.SetRem(call.Args().Num() % call.Args().Denom())
	return nil
}
