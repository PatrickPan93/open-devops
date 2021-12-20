package rpc

import (
	"bufio"
	"io"
	"net"
	"net/rpc"
	"reflect"
	"time"

	"github.com/ugorji/go/codec"

	"github.com/pkg/errors"

	"github.com/toolkits/pkg/net/gobrpc"
)

type RpcCli struct {
	Cli        *gobrpc.RPCClient
	ServerAddr string
}

func InitRpcCli(serverAddr string) *RpcCli {
	return &RpcCli{
		Cli:        nil,
		ServerAddr: serverAddr,
	}
}

// GetCli 如果cli存在就返回 否则new一个
func (r *RpcCli) GetCli() error {
	if r.Cli != nil {
		return nil
	}
	conn, err := net.DialTimeout("tcp", r.ServerAddr, 5*time.Second)
	if err != nil {
		return errors.Wrap(err, "rpc.GetCli: Error while getting rpc cli from server")
	}
	// new a Msg handler
	var mh codec.MsgpackHandle
	// use reflect to define type
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

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
	// 构造RPC Client
	rpcCodec := codec.MsgpackSpecRpc.ClientCodec(bufConn, &mh)
	client := rpc.NewClientWithCodec(rpcCodec)
	r.Cli = gobrpc.NewRPCClient(r.ServerAddr, client, 5*time.Second)
	return nil
}
