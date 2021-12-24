package mem_index

import (
	"log"
	"open-devops/src/common"
	"open-devops/src/modules/server/config"
	"strings"

	ii "github.com/ning1875/inverted-index"
)

type ResourceIndexer interface {
	FlushIndex()                         // 刷新索引方法
	GetIndexReader() *ii.HeadIndexReader // 获取索引reader
	//GetLogger() log.Logger
}

var indexContainer = make(map[string]ResourceIndexer)

func iRegister(name string, ri ResourceIndexer) {
	indexContainer[name] = ri
}

func JudgeResourceIndexExists(name string) bool {
	_, ok := indexContainer[name]
	return ok
}

func Init(ims []*config.IndexModuleConf) {
	loadNum := 0
	loadResource := make([]string, 0)
	for _, i := range ims {
		if !i.Enable {
			log.Printf("mem-index.Init: resouce %s is disabled\n", i.ResourceName)
			continue
		}
		log.Printf("mem-index.Init: resouce %s is enabled\n", i.ResourceName)
		loadNum++
		loadResource = append(loadResource, i.ResourceName)
		switch i.ResourceName {
		case common.ResourceHost:
			mi := &HostIndex{
				Ir:      ii.NewHeadReader(),
				Logger:  nil,
				Modulus: i.Modulus,
				Num:     i.Num,
			}
			iRegister(i.ResourceName, mi)
			//TODO case other resources
		default:
			log.Printf("mem-index.Init: resource %s is not supported now", i.ResourceName)
		}
	}
	log.Printf("mem-index.Init: loadNum %d, details %s", loadNum, strings.Join(loadResource, ","))
}
