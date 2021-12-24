package mem_index

import (
	"context"
	"log"
	"open-devops/src/common"
	"open-devops/src/modules/server/config"
	"strings"
	"sync"
	"time"

	"github.com/ning1875/inverted-index/index"

	"github.com/pkg/errors"

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
				Ir: ii.NewHeadReader(),
				//Logger:  nil,
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

func GetMatchIdsByIndex(req common.ResourceQueryReq) (matchIds []uint64) {
	// 尝试从接口容器中根据查询资源类型获取indexer
	ri, ok := indexContainer[req.ResourceType]
	if !ok {
		log.Printf("mem-index.GetMatchIdsByIndex: target index %s doesn't exist", req.ResourceType)
		return
	}

	// 根据请求中的labels兑换获取matchers
	matchers := common.FormatLabelMatcher(req.Labels)

	// 根据matchers获取对应的IndexReader
	p, err := ii.PostingsForMatchers(ri.GetIndexReader(), matchers...)

	if err != nil {
		log.Printf("%+v", errors.Wrap(err, "mem-index.GetMatchIdsByIndex: Error while getting IndexReader"))
		return
	}
	matchIds, err = index.ExpandPostings(p)
	if err != nil {
		log.Printf("%+v", errors.Wrap(err, "mem-index.GetMatchIdsByIndex: Error while getting matchIds"))
		return
	}
	return
}

func RevertedIndexSyncManager(ctx context.Context) error {
	log.Println("RevertedIndexSyncManager started")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("mem-index.RevertedIndexSyncManager: received term signal.. would be exit soon")
			return nil
		case <-ticker.C:
			log.Printf("sync.RevertedIndexSyncManager: start doIndexFlush %d", len(indexContainer))
			doIndexFlush()
		}
	}
}

func doIndexFlush() {
	var wg sync.WaitGroup
	wg.Add(len(indexContainer))
	for _, ir := range indexContainer {
		go func() {
			defer wg.Done()
			ir.FlushIndex()
		}()
	}
	wg.Wait()
}
