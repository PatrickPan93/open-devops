package rpc

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/rpc"
	"reflect"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

// Server serverType
type Server int

func Start(rpcAddr string) error {
	// new rpc server
	server := rpc.NewServer()
	// register server object to server
	err := server.Register(new(Server))
	/*
		if err != nil {
			return errors.Wrap(err, "Error while registering server to rpc")
		}

	*/

	l, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return errors.Wrap(err, "Start: failed starting rpc server")
	}
	log.Printf("rpc server available: %s\n", rpcAddr)

	// new a Msg handler
	var mh codec.MsgpackHandle
	// use reflect to define type
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
	for {
		// get conn from rpc server
		conn, err := l.Accept()
		// handle it
		if err != nil {
			log.Printf("%+v", errors.Wrap(err, "Error while getting conn from rpc server:"))
			continue
		}
		// 用bufio作为io解析提速
		var bufConn = struct {
			io.Closer
			*bufio.Reader
			*bufio.Writer
		}{
			conn,
			bufio.NewReader(conn),
			bufio.NewWriter(conn),
		}
		// 使用mh handler 来进行解析, 使用bufConn作为提速
		go server.ServeCodec(codec.MsgpackSpecRpc.ServerCodec(bufConn, &mh))
	}
}
