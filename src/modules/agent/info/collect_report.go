package info

import (
	"context"
	"log"
	"open-devops/src/common"
	"open-devops/src/models"
	"open-devops/src/modules/agent/rpc"
	"time"

	"github.com/pkg/errors"
)

const (
	TencentCloudGettingLocalIP = `curl http://metadata.tencentyun.com/meta-data/local-ipv4`
	GettingMachineSNGYY        = `curl -s http://169.254.169.254/a/meta-data/instance-id`
	GettingMachineSNLocal      = `dmidecode -s system-serial-number | tail -n 1 | tr -d "\n"`

	GettingCPUCoresShell = `cat /proc/cpuinfo | grep processor | wc -l | tr -d "\n"`
	GettingMemTotal      = `cat /proc/meminfo | grep MemTotal | awk '{printf "%d",$2/1024/1024}'`

	GettingDiskTotal = `df -m | grep '/dev/' | grep -v /var/lib | grep -v tmpfs | awk '{sum+=$2};END{printf "%d",sum/1024}'`
)

// TickerInfoCollectAndReport Using time ticker to run CollectBaseInfo
func TickerInfoCollectAndReport(cli *rpc.RpcCli, ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	log.Println("info.TickerInfoCollectAndReport started")
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("info.TickerInfoCollectAndReport: ticker would be stopped because of term signal")
			return nil
		case <-ticker.C:
			CollectBaseInfo(cli)
		}
	}
}

// CollectBaseInfo Getting Machine Info
func CollectBaseInfo(cli *rpc.RpcCli) {
	sn, err := common.ShellCommand(GettingMachineSNGYY)
	if err != nil || sn == "" {
		sn, err = common.ShellCommand(GettingMachineSNLocal)
		if err != nil {
			log.Printf("%+v\n", errors.Wrap(err, "info.CollectBaseInfo: Error while getting Machine SN by GYY and Local way"))
			return
		}
	}
	log.Printf("info.CollectBaseInfo: sucessfully getting Machine SN: %s", sn)

	cpuCore, err := common.ShellCommand(GettingCPUCoresShell)
	if err != nil {
		log.Printf("%+v\n", errors.Wrap(err, "info.CollectBaseInfo: Error while getting Machine CPU Core"))
		return
	}
	log.Printf("info.CollectBaseInfo: sucessfully getting Machine CPU Core: %s", cpuCore)

	memTotal, err := common.ShellCommand(GettingMemTotal)
	if err != nil {
		log.Printf("%+v\n", errors.Wrap(err, "info.CollectBaseInfo: Error while getting Machine Mem Total"))
		return
	}
	log.Printf("info.CollectBaseInfo: sucessfully getting Machine Mem Total: %s", memTotal)

	diskTotal, err := common.ShellCommand(GettingDiskTotal)
	if err != nil {
		log.Printf("%+v\n", errors.Wrap(err, "info.CollectBaseInfo: Error while getting Machine Disk Total"))
		return
	}
	log.Printf("info.CollectBaseInfo: sucessfully getting Machine Disk Total: %s", diskTotal)

	hostOjb := &models.AgentCollectInfo{
		SN:       sn,
		CPU:      cpuCore,
		Mem:      memTotal,
		Disk:     diskTotal,
		IPAddr:   common.GetLocalIP(),
		Hostname: common.GetHostname(),
	}
	err = cli.GetCli()
	if err != nil {
		log.Printf("%+v", err)
	}
	err = cli.HostInfoReport(hostOjb)
	if err != nil {
		log.Printf("%+v", err)
	}
}
