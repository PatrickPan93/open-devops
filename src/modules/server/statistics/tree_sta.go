package statistics

import (
	"context"
	"log"
	"open-devops/src/common"
	mem_index "open-devops/src/modules/server/mem-index"
	"open-devops/src/modules/server/metric"
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
	for resourceType, ir := range irs {
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

	}
}
