package rpc

import (
	"fmt"
	"log"
	"open-devops/src/models"

	"github.com/pkg/errors"
)

func (r *RpcCli) HostInfoReport(info *models.AgentCollectInfo) error {
	var msg string
	err := r.GetCli()
	if err != nil {
		log.Printf("%+v", err)
		return errors.Wrap(err, "rpc.Ping: Error while getting cli to ping server")
	}
	// calling rpc server
	// Method为server端注册的别名(type Server int)
	// 后两个参数对应是server Ping fun的input和output
	err = r.Cli.Call("Server.HostInfoReport", info, &msg)
	if err != nil {
		//r.Cli.Close()
		log.Println(fmt.Sprintf("rpc.HostInfoReport: Error while reporting info to server: %+v", err))
		return errors.Wrap(err, "rpc.HostInfoReport: Error while reporting info to server")
	}
	log.Printf("rpc.Ping: reporting info to server successfully: serverAddr: %s, Msg: %s", r.ServerAddr, msg)
	return nil
}
