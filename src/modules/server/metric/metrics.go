package metric

import (
	"open-devops/src/common"

	"github.com/prometheus/client_golang/prometheus"
)

// 自定义prometheus采集指标
var (
	// Index
	IndexFlushDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_index_flush_last_seconds",
		Help: "Duration of index flush",
	}, []string{common.LABEL_RESOURCE_TYPE})

	IndexResourceNumCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_num_count",
		Help: "Number of index resources",
	}, []string{common.LABEL_RESOURCE_TYPE})

	// Public Cloud
	PublicCloudSyncDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "public_cloud_sync_duration",
		Help: "Duration of public cloud sync",
	}, []string{common.LABEL_RESOURCE_TYPE})

	PublicCloudSyncResourceNumCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "public_cloud_resource_num_count",
		Help: "Number of public cloud resources",
	}, []string{common.LABEL_RESOURCE_TYPE})

	ResourceNumRegionCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_num_region_count",
		Help: "Number of resources group by region tag",
	}, []string{common.LABEL_RESOURCE_TYPE, common.LABEL_REGION})

	ResourceNumCloudProviderCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_num_cloud_provider_count",
		Help: "Number of resources group by cloud provider tag",
	}, []string{common.LABEL_RESOURCE_TYPE, common.LABEL_CLOUD_PROVIDER})
)

// NewMetrics 注册Metrics
func NewMetrics() {
	// Index
	prometheus.DefaultRegisterer.MustRegister(IndexFlushDuration)
	prometheus.DefaultRegisterer.MustRegister(IndexResourceNumCount)
	// Public Cloud
	prometheus.DefaultRegisterer.MustRegister(PublicCloudSyncDuration)
	prometheus.DefaultRegisterer.MustRegister(PublicCloudSyncResourceNumCount)

	prometheus.DefaultRegisterer.MustRegister(ResourceNumRegionCount)
	prometheus.DefaultRegisterer.MustRegister(ResourceNumCloudProviderCount)
}