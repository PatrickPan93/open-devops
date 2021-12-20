package rpc

import "fmt"

// Ping First rpc func
func (*Server) Ping(input string, output *string) error {
	fmt.Println(input)
	*output = "copy"
	return nil
}
