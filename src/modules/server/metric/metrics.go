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

	GPAAllNumCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_num_count",
		Help: "Num gpa of all",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE})

	// GPA
	GPAAllNumRegionCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_region_num_count",
		Help: "Num gpa of all with tag region",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE, common.LABEL_REGION})

	GPAAllNumCloudProviderCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_cloud_provider_num_count",
		Help: "Num gpa of all with tag cloud provider",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE, common.LABEL_CLOUD_PROVIDER})

	GPAAllNumClusterCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_cluster_num_count",
		Help: "Num gpa of all with tag cluster",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE, common.LABEL_CLUSTER})

	GPAAllNumInstanceTypeCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpa_all_instance_type_num_count",
		Help: "Num gpa of all with tag instance type",
	}, []string{common.LABEL_GPA_NAME, common.LABEL_RESOURCE_TYPE, common.LABEL_INSTANCE_TYPE})

	// Host 特殊
	GPAHostCpuCores = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpb_host_cpu_cores",
		Help: "Num gpa cpu cores of ecs",
	}, []string{common.LABEL_GPA_NAME})

	GPAHostMemGbs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpb_host_mem_gbs",
		Help: "Num gpa mem gbs of ecs",
	}, []string{common.LABEL_GPA_NAME})

	GPAHostDiskGbs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gpb_host_disk_gbs",
		Help: "Num gpa disk gbs of ecs",
	}, []string{common.LABEL_GPA_NAME})
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

	prometheus.DefaultRegisterer.MustRegister(GPAAllNumCount)

	prometheus.DefaultRegisterer.MustRegister(GPAAllNumClusterCount)
	prometheus.DefaultRegisterer.MustRegister(GPAAllNumInstanceTypeCount)
	prometheus.DefaultRegisterer.MustRegister(GPAAllNumRegionCount)
	prometheus.DefaultRegisterer.MustRegister(GPAAllNumCloudProviderCount)

	prometheus.DefaultRegisterer.MustRegister(GPAHostCpuCores)
	prometheus.DefaultRegisterer.MustRegister(GPAHostMemGbs)
	prometheus.DefaultRegisterer.MustRegister(GPAHostDiskGbs)

}
