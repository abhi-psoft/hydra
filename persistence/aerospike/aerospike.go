package aerospike

import (
	"fmt"

	"github.com/aerospike/aerospike-client-go/v7"
)

// if true {
// 	client, err := aerospike.NewClient("127.0.0.1", 3000)
// 	if err != nil {
// 		_, _ = fmt.Println("Failed to connect to Aerospike: ", err)
// 	}
// 	defer client.Close()

// 	writer := aero.NewAerospikeWriter(client)

// 	writer.CreateClient(r.Context(), &c)
// }

var Client *aerospike.Client

func NewAerospikeConn() error {
	var err error
	fmt.Println("NewAerospikeConn Start")
	Client, err = aerospike.NewClient("127.0.0.1", 3000)
	if err != nil {
		_, _ = fmt.Println("Failed to connect to Aerospike: ", err)
		return err
	}
	defer Client.Close()

	return nil
}
