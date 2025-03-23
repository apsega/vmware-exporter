package collector

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

var (
	clusterLabels       []string = []string{"cluster_name", "cluster_id"}
	clusterCpuEffective          = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "cluster",
			Name:      "cpu_effective_mhz",
			Help:      "Effective Cluster CPU in MHz",
		},
		clusterLabels,
	)
	clusterCpuNum = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "cluster",
			Name:      "cpu_cores_total",
			Help:      "Total Cluster CPU cores number",
		},
		clusterLabels,
	)
	clusterCpuTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "cluster",
			Name:      "cpu_mhz_total",
			Help:      "Total Cluster CPU in MHz",
		},
		clusterLabels,
	)
	clusterHostsEffective = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "cluster",
			Name:      "hosts_effective_total",
			Help:      "Effective Cluster hosts",
		},
		clusterLabels,
	)
	clusterHostsNum = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "cluster",
			Name:      "hosts_total",
			Help:      "Total Cluster hosts",
		},
		clusterLabels,
	)
	clusterMemoryEffective = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "cluster",
			Name:      "memory_effective_bytes",
			Help:      "Effective Cluster memory in bytes",
		},
		clusterLabels,
	)
	clusterMemoryTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "cluster",
			Name:      "memory_bytes_total",
			Help:      "Total Cluster memory in bytes",
		},
		clusterLabels,
	)
	clusterThreadsNum = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "cluster",
			Name:      "threads_total",
			Help:      "Total Cluster threads",
		},
		clusterLabels,
	)
)

func ExportClusterMetrics(ctx context.Context, client *govmomi.Client) error {

	c := client.Client
	// Create a view manager
	m := view.NewManager(c)

	// Create a container view for clusters
	containerView, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"ClusterComputeResource"}, true)
	if err != nil {
		fmt.Printf("Error creating container view: %v\n", err)
		return err
	}
	defer containerView.Destroy(ctx)

	// Retrieve a list of clusters
	var clusters []mo.ClusterComputeResource
	if err = containerView.Retrieve(ctx, []string{"ClusterComputeResource"}, nil, &clusters); err != nil {
		fmt.Printf("Error retrieving clusters: %v\n", err)
		return err
	}

	// Iterate through the clusters and print provisioning vs. available resources
	for _, cluster := range clusters {

		clusterName := cluster.Name
		clusterID := cluster.Self.Reference().Value

		for _, host := range cluster.Host {
			populateHostMapping(host.Value, clusterName)
		}

		labels := prometheus.Labels{
			"cluster_name": clusterName,
			"cluster_id":   clusterID,
		}
		clusterCpuEffective.With(labels).Set(
			float64(cluster.Summary.GetComputeResourceSummary().EffectiveCpu),
		)
		clusterCpuNum.With(labels).Set(
			float64(cluster.Summary.GetComputeResourceSummary().NumCpuCores),
		)
		clusterCpuTotal.With(labels).Set(
			float64(cluster.Summary.GetComputeResourceSummary().TotalCpu),
		)
		clusterHostsEffective.With(labels).Set(
			float64(cluster.Summary.GetComputeResourceSummary().NumEffectiveHosts),
		)
		clusterHostsNum.With(labels).Set(
			float64(cluster.Summary.GetComputeResourceSummary().NumHosts),
		)
		clusterMemoryEffective.With(labels).Set(
			float64(cluster.Summary.GetComputeResourceSummary().EffectiveMemory),
		)
		clusterMemoryTotal.With(labels).Set(
			float64(cluster.Summary.GetComputeResourceSummary().TotalMemory),
		)
		clusterThreadsNum.With(labels).Set(
			float64(cluster.Summary.GetComputeResourceSummary().NumCpuThreads),
		)
	}
	return nil
}
