package cloud_sync

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"open-devops/src/common"
	"open-devops/src/models"
	"time"
)

type HostSync struct {
	CloudAlibaba
	CloudTencent
	TableName string
}

func (hs *HostSync) sync() {
	start := time.Now()
	// use mock data as the sync data for testing
	resourceHosts := genMockResourceHost()
	// get uuid hash map from db
	uuidHashM, err := models.GetHostUidAndHash()
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	// - toAddSet, toModSet存放的都是host对象
	toAddSet := make([]models.ResourceHost, 0)
	toModSet := make([]models.ResourceHost, 0)
	toDelSet := make([]string, 0)

	var toAddNum, toModNum, toDelNum int
	var suAddNum, suModNum, suDelNum int
	localUidSet := make(map[string]struct{})

	for _, h := range resourceHosts {
		localUidSet[h.Uid] = struct{}{}
		hash, ok := uuidHashM[h.Uid]
		if !ok {
			// 说明本地不存在该uid, 但是公有云存在,则新增
			toAddSet = append(toAddSet, h)
			toAddNum++
			continue
		}
		// 存在则判断hash是否相同
		if hash == h.Hash {
			// 相同说明不需要进行同步
			continue
		}
		//说明Uid相通但是hash不同,某些字段发生变更
		toModSet = append(toModSet, h)
		toModNum++
	}
	// 获取出db的uid与公有云uid比较
	for uid := range uuidHashM {

		// 如果本地db的uid不存在于公有云结果中，说明需要被删除
		if _, ok := localUidSet[uid]; !ok {
			toDelSet = append(toDelSet, uid)
			toDelNum++
		}
	}
	// 上述的判断流程
	// 下面进行具体执行操作
	for _, h := range toAddSet {
		_, err := h.AddOne()
		if err != nil {
			log.Printf("%+v\n", err)
			continue
		}
		suAddNum++
	}
	for _, h := range toModSet {
		_, err := h.UpdateByUid(h.Uid)
		if err != nil {
			log.Printf("%+v\n", err)
			continue
		}
		suModNum++
	}
	// 删除
	if len(toDelSet) > 0 {
		num, err := models.BatchDeleteResource(common.ResourceHost, "uid", toDelSet)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		suDelNum = int(num)
	}
	timeTook := time.Since(start)
	log.Printf("cloud_sync.sync:\n public.cloud.num: %d\n, db.num: %d\n, toAddNum, %d\n, toModNum: %d\n, toDelNum: %d\n, suAddNum: %d\n, suModNum: %d\n, suDelNum: %d\n, timeTook: %s\n", len(resourceHosts), len(uuidHashM), toAddNum, toModNum, toDelNum, suAddNum, suModNum, suDelNum, timeTook)
}

// 模拟公有云资源
func genMockResourceHost() []models.ResourceHost {
	rand.Seed(time.Now().UnixNano())
	// Rand g.p.a
	randGs := []string{"inf", "ads", "web", "sys"}
	randPs := []string{"monitor", "cicd", "k8s", "mq"}
	randAs := []string{"kafka", "prometheus", "zookeeper", "elasticsearch"}

	// rand hardware
	randCpus := []string{"4", "8", "16", "32", "64", "128", "256"}
	randMems := []string{"16", "32", "64", "128", "257"}
	randDisk := []string{"50", "100", "150", "200", "300", "500"}

	// 标签tags
	randMapKeys := []string{"arch", "idc", "os", "job"}
	randMapValues := []string{"linux", "beijing", "windows", "arm64", "amd64", "darwin"}

	// 公有云标签
	randRegions := []string{"beijing", "guangzhou", "xinjiang", "ningxia", "shandong", "tianjin"}
	randCloudProviders := []string{"alibaba", "tencent", "aws", "gcp", "huawei", "azure"}
	randCluster := []string{"bigdata", "inf", "middleware", "business"}
	randInst := []string{"4c8g", "4c16g", "8c32g", "16c64g", "32c128g"}

	frn := func(n int) int {
		return rand.Intn(n - 1)
	}
	frNum := func() int {
		return int(rand.Int63n(60-25) + 25)
	}
	hs := make([]models.ResourceHost, 0)

	for i := 0; i < frNum(); i++ {
		randN := i
		name := fmt.Sprintf("genMockResourceHost_host_%d", randN)
		ips := []string{fmt.Sprintf("8.8.8.%d", randN)}
		ipJ, _ := json.Marshal(ips)
		h := models.ResourceHost{
			Name:          name,
			PrivateIps:    ipJ,
			CPU:           randCpus[frn(len(randCpus))],
			Mem:           randMems[frn(len(randMems))],
			Disk:          randDisk[frn(len(randDisk))],
			StreeGroup:    randGs[frn(len(randGs))],
			StreeProduct:  randPs[frn(len(randPs))],
			StreeApp:      randAs[frn(len(randAs))],
			Region:        randRegions[frn(len(randRegions))],
			CloudProvider: randCloudProviders[frn(len(randCloudProviders))],
			InstanceType:  randInst[frn(len(randInst))],
		}
		// gen tags
		tagM := make(map[string]string)
		for _, v := range randMapKeys {
			tagM[v] = randMapValues[frn(len(randMapValues))]
		}
		// cluster
		tagM["cluster"] = randCluster[(frn(len(randCluster)))]
		tagMj, _ := json.Marshal(tagM)
		h.Tags = tagMj
		// gen hash
		hash := h.GetHash()
		h.Hash = hash
		// genUUID
		md5o := md5.New()
		md5o.Write([]byte(name))
		h.Uid = hex.EncodeToString(md5o.Sum(nil))
		hs = append(hs, h)
	}
	return hs
}
