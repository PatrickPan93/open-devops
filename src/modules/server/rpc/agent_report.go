package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"open-devops/src/models"

	"github.com/pkg/errors"
)

// HostInfoReport First rpc func
func (*Server) HostInfoReport(input models.AgentCollectInfo, output *string) error {

	bytes, err := json.Marshal(input)
	if err != nil {
		log.Println(fmt.Sprintf("rpc.HostInfoReport: Error while receiving info from agent: %+v", err))
		return errors.Wrap(err, "Error while translating struct to json")
	}

	*output = "Already received data"
	fmt.Println(string(bytes))
	ips := []string{input.IPAddr}
	ipJ, _ := json.Marshal(ips)
	if input.SN == "" {
		input.SN = input.Hostname
	}
	if input.SN == "" {
		*output = "sn.empty"
		return nil
	}
	// 获取对象Uid
	rh := models.ResourceHost{
		Uid:        input.SN,
		Name:       input.Hostname,
		PrivateIps: ipJ,
		SN:         input.SN,
		CPU:        input.CPU,
		Mem:        input.Mem,
		Disk:       input.Disk,
		Hostname:   input.Hostname,
	}
	hash := rh.GetHash()

	// 用Uid去db中获取之前的结果,再根据两者的Hash是否一致决定更改
	rhUid := models.ResourceHost{Uid: input.SN}
	rhUidDb, err := rhUid.GetOne()
	if err != nil {
		*output = fmt.Sprintf("db_error_%v", err)
		return nil
	}
	if rhUidDb == nil {
		// 说明Uid不存在,需要插入数据
		rh.Hash = hash
		_, err := rh.AddOne()
		if err != nil {
			*output = fmt.Sprintf("db_error_%v", err)
			return nil
		}
		*output = "insert_success"
		return nil
	}

	// 若Uid存在则需要判断Hash
	if rhUidDb.Hash != hash {
		rh.Hash = hash
		rowAffected, err := rh.UpdateByUid(rh.Uid)
		if err != nil {
			*output = fmt.Sprintf("db_error_%v", err)
			return nil
		}
		*output = fmt.Sprintf("update_success, %d rows affected", rowAffected)
		return nil
	}
	*output = "hash equal.. nothing need to do"
	return nil
}
