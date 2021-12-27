package statistics

import (
	"context"
	"log"
	"open-devops/src/common"
	"open-devops/src/models"
	mem_index "open-devops/src/modules/server/mem-index"
	"open-devops/src/modules/server/metric"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TreeNodeStatisticsManager(ctx context.Context) error {
	log.Println("TreeNodeStatisticsManager started")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("statistics.TreeNodeStatisticsManager: received term signal.. would be exit soon")
			return nil
		case <-ticker.C:
			log.Println("statistics.TreeNodeStatisticsManager: start statisticsWork")
			statisticsWork()
		}
	}
}

func statisticsWork() {
	irs := mem_index.GetAllResourceIndexReader()
	log.Printf("statistics.statisticsWork: got %d of indexReader", len(irs))

	// 获取所有g.p.a列表
	qReq := &common.NodeCommonReq{
		QueryType: 5,
	}
	allGPAS := models.StreePathQuery(qReq)

	for resourceType, ir := range irs {
		go func() {

			ir := ir
			// 按照region进行分布
			s := ir.GetIndexReader().GetGroupByLabel(common.LABEL_REGION)
			for _, i := range s.Group {
				// 打点 prometheus
				metric.ResourceNumRegionCount.With(prometheus.Labels{common.LABEL_RESOURCE_TYPE: resourceType, common.LABEL_REGION: i.Name}).Set(float64(i.Value))
			}

			// 按照cloud provider进行分布
			p := ir.GetIndexReader().GetGroupByLabel(common.LABEL_CLOUD_PROVIDER)
			for _, i := range p.Group {
				// 打点 prometheus
				metric.ResourceNumCloudProviderCount.With(prometheus.Labels{common.LABEL_RESOURCE_TYPE: resourceType, common.LABEL_CLOUD_PROVIDER: i.Name}).Set(float64(i.Value))
			}

			// GPA 维度
			// 拿出gpa串
			for _, gpa := range allGPAS {
				ss := strings.Split(gpa, ".")
				if len(ss) != 3 {
					continue
				}
				// split后构造matcher
				g, p, a := ss[0], ss[1], ss[2]
				csG := &common.SingleTagReq{
					Key:   common.LABEL_STREE_G,
					Value: g,
					Type:  1,
				}
				csP := &common.SingleTagReq{
					Key:   common.LABEL_STREE_P,
					Value: p,
					Type:  1,
				}
				csA := &common.SingleTagReq{
					Key:   common.LABEL_STREE_A,
					Value: a,
					Type:  1,
				}
				matcherG := []*common.SingleTagReq{
					csG,
				}
				matcherGP := []*common.SingleTagReq{
					csG, csP,
				}
				matcherGPA := []*common.SingleTagReq{
					csG, csP, csA,
				}

				// g.p.a 各级资源
				gpaNumWork(resourceType, g, matcherG, metric.GPAAllNumCount)
				gpaNumWork(resourceType, g+"."+p, matcherGP, metric.GPAAllNumCount)
				gpaNumWork(resourceType, g+"."+p+"."+a, matcherGPA, metric.GPAAllNumCount)
				// 按照g来进行统计打点(group by g)
				// g
				gpaLabelNumWork(resourceType, common.LABEL_REGION, g, matcherG, ir, metric.GPAAllNumRegionCount)
				gpaLabelNumWork(resourceType, common.LABEL_CLOUD_PROVIDER, g, matcherG, ir, metric.GPAAllNumCloudProviderCount)
				gpaLabelNumWork(resourceType, common.LABEL_CLUSTER, g, matcherG, ir, metric.GPAAllNumClusterCount)
				gpaLabelNumWork(resourceType, common.LABEL_INSTANCE_TYPE, g, matcherG, ir, metric.GPAAllNumInstanceTypeCount)

				// p
				gpaLabelNumWork(resourceType, common.LABEL_REGION, g+"."+p, matcherGP, ir, metric.GPAAllNumRegionCount)
				gpaLabelNumWork(resourceType, common.LABEL_CLOUD_PROVIDER, g+"."+p, matcherGP, ir, metric.GPAAllNumCloudProviderCount)
				gpaLabelNumWork(resourceType, common.LABEL_CLUSTER, g+"."+p, matcherGP, ir, metric.GPAAllNumClusterCount)
				gpaLabelNumWork(resourceType, common.LABEL_INSTANCE_TYPE, g+"."+p, matcherGP, ir, metric.GPAAllNumInstanceTypeCount)

				// a
				gpaLabelNumWork(resourceType, common.LABEL_REGION, g+"."+p+"."+a, matcherGPA, ir, metric.GPAAllNumRegionCount)
				gpaLabelNumWork(resourceType, common.LABEL_CLOUD_PROVIDER, g+"."+p+"."+a, matcherGPA, ir, metric.GPAAllNumCloudProviderCount)
				gpaLabelNumWork(resourceType, common.LABEL_CLUSTER, g+"."+p+"."+a, matcherGPA, ir, metric.GPAAllNumClusterCount)
				gpaLabelNumWork(resourceType, common.LABEL_INSTANCE_TYPE, g+"."+p+"."+a, matcherGPA, ir, metric.GPAAllNumInstanceTypeCount)

				if resourceType == common.ResourceHost {
					// g
					hostSpecial(resourceType, common.LABEL_CPU, g, matcherG, ir, metric.GPAHostCpuCores)
					hostSpecial(resourceType, common.LABEL_MEM, g, matcherG, ir, metric.GPAHostMemGbs)
					hostSpecial(resourceType, common.LABEL_DISK, g, matcherG, ir, metric.GPAHostDiskGbs)
					// g.p
					hostSpecial(resourceType, common.LABEL_CPU, g+"."+p, matcherGP, ir, metric.GPAHostCpuCores)
					hostSpecial(resourceType, common.LABEL_MEM, g+"."+p, matcherGP, ir, metric.GPAHostMemGbs)
					hostSpecial(resourceType, common.LABEL_DISK, g+"."+p, matcherGP, ir, metric.GPAHostDiskGbs)
					// g.p.a
					hostSpecial(resourceType, common.LABEL_CPU, g+"."+p+"."+a, matcherGPA, ir, metric.GPAHostCpuCores)
					hostSpecial(resourceType, common.LABEL_MEM, g+"."+p+"."+a, matcherGPA, ir, metric.GPAHostMemGbs)
					hostSpecial(resourceType, common.LABEL_DISK, g+"."+p+"."+a, matcherGPA, ir, metric.GPAHostDiskGbs)
				}
			}
		}()
	}
}

// 封装gpa通用方法
func gpaLabelNumWork(resourceType string, targetLabel string, gpaName string, matcher []*common.SingleTagReq, ir mem_index.ResourceIndexer, ms *prometheus.GaugeVec) {
	req := common.ResourceQueryReq{
		ResourceType: resourceType,
		Labels:       matcher,
		TargetLabel:  targetLabel,
	}
	matchIds := mem_index.GetMatchIdsByIndex(req)

	statsRs := ir.GetIndexReader().GetGroupDistributionByLabel(req.TargetLabel, matchIds)

	for _, x := range statsRs.Group {
		ms.With(prometheus.Labels{
			common.LABEL_GPA_NAME:      gpaName,
			common.LABEL_RESOURCE_TYPE: resourceType, targetLabel: x.Name}).Set(float64(x.Value))
	}
}

// 通过索引的getGroupByLabel获取个数分布
// 每个g.p.a在每种资源上的统计
func gpaNumWork(resourceType string, gpaName string, matcher []*common.SingleTagReq, ms *prometheus.GaugeVec) {
	req := common.ResourceQueryReq{
		ResourceType: resourceType,
		Labels:       matcher,
	}
	matchIds := mem_index.GetMatchIdsByIndex(req)
	if len(matchIds) > 0 {
		ms.With(prometheus.Labels{
			common.LABEL_GPA_NAME:      gpaName,
			common.LABEL_RESOURCE_TYPE: resourceType,
		}).Set(float64(len(matchIds)))
	}
}

// host特殊的
func hostSpecial(resourceType string, targetLabel string, gpaName string, matcher []*common.SingleTagReq, ir mem_index.ResourceIndexer, ms *prometheus.GaugeVec) {
	req := common.ResourceQueryReq{
		ResourceType: resourceType,
		Labels:       matcher,
		TargetLabel:  targetLabel,
	}
	matchIds := mem_index.GetMatchIdsByIndex(req)

	statsRe := ir.GetIndexReader().GetGroupDistributionByLabel(targetLabel, matchIds)

	var all uint64
	for _, x := range statsRe.Group {
		num, _ := strconv.Atoi(x.Name)
		all += uint64(num) * x.Value
	}
	if all > 0 {
		ms.With(prometheus.Labels{common.LABEL_GPA_NAME: gpaName}).Set(float64(all))
	}

}
