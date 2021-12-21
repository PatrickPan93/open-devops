package rpc

import (
	"log"

	"github.com/pkg/errors"
)

func (r *RpcCli) Ping() error {
	var msg string
	err := r.GetCli()
	if err != nil {
		return errors.Wrap(err, "rpc.Ping: Error while getting cli to ping server")
	}
	// calling rpc server
	// Method为server端注册的别名(type Server int)
	// 后两个参数对应是server Ping fun
	err = r.Cli.Call("Server.Ping", "agent01", &msg)
	if err != nil {
		return errors.Wrap(err, "rpc.Ping: Error while pinging server")
	}
	log.Printf("rpc.Ping: ping server successfully: serverAddr: %s, Msg: %s", r.ServerAddr, msg)
	return nil
}
