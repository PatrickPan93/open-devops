package info

import (
	"context"
	"log"
	"open-devops/src/common"
	"time"

	"github.com/pkg/errors"
)

const (
	TencentCloudGettingLocalIP = `curl http://metadata.tencentyun.com/meta-data/local-ipv4`
	GettingMachineSNGYY        = `curl -s http://169.254.169.254/a/meta-data/instance-id`
	GettingMachineSNLocal      = `dmidecode -s system-serial-number | tail -n 1`

	GettingCPUCoresShell = `cat /proc/cpuinfo | grep processor | wc -l`
	GettingMemTotal      = `cat /proc/meminfo | grep MemTotal | awk '{printf "%d",$2/1024/1024}'`
)

// TickerInfoCollectAndReport Using time ticker to run CollectBaseInfo
func TickerInfoCollectAndReport(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	log.Println("info.TickerInfoCollectAndReport started")
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("info.TickerInfoCollectAndReport: ticker would be stopped because of term signal")
			return nil
		case <-ticker.C:
			CollectBaseInfo()
		}
	}
}

// CollectBaseInfo Getting Machine SN
func CollectBaseInfo() {
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
		log.Printf("%+v\n", errors.Wrap(err, "info.CollectBaseInfo: Error while getting Machine CPU Number"))
		return
	}
	log.Printf("info.CollectBaseInfo: sucessfully getting Machine CPU Num: %s", cpuCore)

	memTotal, err := common.ShellCommand(GettingMemTotal)
	if err != nil {
		log.Printf("%+v\n", errors.Wrap(err, "info.CollectBaseInfo: Error while getting Machine Mem Total"))
		return
	}
	log.Printf("info.CollectBaseInfo: sucessfully getting Machine Mem Total: %s", memTotal)
}
