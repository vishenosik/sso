package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

func main() {

	client := newClient()

	err := setup(client)
	if err != nil {
		log.Fatal(fmt.Errorf("main %w", err))
	}

	txn := client.NewTxn()
	defer txn.Discard(context.TODO())
	q(txn)
}

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	dialOpts := append([]grpc.DialOption{},
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	d, err := grpc.Dial("localhost:9080", dialOpts...)

	if err != nil {
		log.Fatal(fmt.Errorf("newClient %w", err))
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func setup(c *dgo.Dgraph) error {
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	return c.Alter(context.Background(), &api.Operation{
		Schema: `
			name: string @index(term) .
			balance: int .
		`,
	})
}

func q(txn *dgo.Txn) {
	// Query the balance for Alice and Bob.
	const q = `
{
  people(func: has(follows)) {
    name
    age
  }
}

	`
	resp, err := txn.Query(context.Background(), q)
	if err != nil {
		log.Fatal(fmt.Errorf("q-1 %w", err))
	}

	// After we get the balances, we have to decode them into structs so that
	// we can manipulate the data.
	var decode struct {
		People []struct {
			Name string
			Age  int
		}
	}
	if err := json.Unmarshal(resp.GetJson(), &decode); err != nil {
		log.Fatal(fmt.Errorf("q-2 %w", err))
	}

	fmt.Println(decode, string(resp.GetJson()))

}
