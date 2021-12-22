package cloud_sync

import (
	"context"
	"log"
	"open-devops/src/common"
	"sync"
	"time"
)

type CloudResource interface {
	sync()
}

type CloudAlibaba struct {
}

type CloudTencent struct {
}

// 接口容器,承载多个资源接口的同步
var (
	cloudResourceContainer = make(map[string]CloudResource)
)

// 注册函数
func cRegister(name string, cr CloudResource) {
	cloudResourceContainer[name] = cr
}

func init() {
	hs := &HostSync{
		TableName: common.ResourceHost,
	}
	cRegister(common.ResourceHost, hs)
}

// CloudSyncManager 管理接口容器管理端
func CloudSyncManager(ctx context.Context) error {
	log.Println("CloudSyncManager started")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("sync.CloudSyncManager: received term signal.. would be exit soon")
			return nil
		case <-ticker.C:
			log.Printf("sync.CloudSyncManager: start doCloudSync resource_num %d\n", len(cloudResourceContainer))
			doCloudSync()
		}
	}
}

// 用wg对任务进行并发管理
func doCloudSync() {
	var wg sync.WaitGroup
	wg.Add(len(cloudResourceContainer))
	for _, sy := range cloudResourceContainer {
		go func() {
			defer wg.Done()
			sy.sync()
		}()
	}
	wg.Wait()
}
