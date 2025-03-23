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
	hostLabels    []string = []string{"host_name", "host_id", "datacenter", "cluster_name"}
	hostAvailPMem          = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "available_pmem_bytes",
			Help:      "Host available persistent memory in bytes",
		},
		hostLabels,
	)
	hostCpuAllocRes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_allocation_reservation_mhz",
			Help:      "Host CPU allocation reservation in Mhz",
		},
		hostLabels,
	)
	hostCpuAllocLim = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_allocation_limit_mhz",
			Help:      "Host CPU allocation limit in Mhz",
		},
		hostLabels,
	)
	hostCpuAllocOver = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_allocation_overhead_mhz",
			Help:      "Host CPU allocation overhead in Mhz",
		},
		hostLabels,
	)
	hostCpuCores = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_cores_total",
			Help:      "Total Host CPU cores number",
		},
		hostLabels,
	)
	hostCpuFree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_free_mhz",
			Help:      "Free Host CPU in Mhz",
		},
		hostLabels,
	)
	hostCpuOverallUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_usage_mhz",
			Help:      "Overall Host CPU usage in Mhz",
		},
		hostLabels,
	)
	hostCpuMhz = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_core_mhz",
			Help:      "Host CPU core Mhz",
		},
		hostLabels,
	)
	hostCpuTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_mhz_total",
			Help:      "Total Host CPU in Mhz",
		},
		hostLabels,
	)
	hostCpuThreads = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "cpu_threads_total",
			Help:      "Total Host CPU threads number",
		},
		hostLabels,
	)
	hostMemoryAllocRes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "memory_allocation_bytes",
			Help:      "Host memory allocation in bytes",
		},
		hostLabels,
	)
	hostMemoryAllocLim = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "memory_allocation_limit_bytes",
			Help:      "Host memory allocation limit in bytes",
		},
		hostLabels,
	)
	hostMemoryFree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "memory_free_bytes",
			Help:      "Host memory free in bytes",
		},
		hostLabels,
	)
	hostMemorySize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "memory_bytes_total",
			Help:      "Total Host memory in bytes",
		},
		hostLabels,
	)
	hostMemoryOverallUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "memory_usage_bytes",
			Help:      "Overall Host memory usage in bytes",
		},
		hostLabels,
	)
	hostNicsNum = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "nics_total",
			Help:      "Total Host NICs number",
		},
		hostLabels,
	)
	hostUptime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "host",
			Name:      "uptime_seconds",
			Help:      "Host uptime in seconds",
		},
		hostLabels,
	)
)

func ExportHostMetrics(ctx context.Context, c *govmomi.Client) error {
	// Create a view manager
	m := view.NewManager(c.Client)

	// Create a container view for the datacenters
	containerView, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		fmt.Printf("Error creating container view: %v\n", err)
	}
	defer containerView.Destroy(ctx)

	// Retrieve a list of datacenters
	var datacenters []mo.Datacenter
	err = containerView.Retrieve(ctx, []string{"Datacenter"}, nil, &datacenters)
	if err != nil {
		fmt.Printf("Error retrieving datacenters: %v\n", err)
	}

	// Iterate through the datacenters and list hosts
	for _, dc := range datacenters {

		// Create a container view for the hosts in the datacenter
		containerView, err := m.CreateContainerView(ctx, dc.Reference(), []string{"HostSystem"}, true)
		if err != nil {
			fmt.Printf("Error creating container view for hosts: %v\n", err)
		}
		defer containerView.Destroy(ctx)

		// Retrieve a list of hosts in the datacenter
		var hosts []mo.HostSystem

		err = containerView.Retrieve(ctx, []string{"HostSystem"}, nil, &hosts)
		if err != nil {
			fmt.Printf("Error retrieving hosts: %v\n", err)
		}

		for _, host := range hosts {

			hostID := host.Self.Reference().Value
			populateVirtualMachineMapping(hostID, host.Name)
			populateHostConfig(host.Name, float64(host.Summary.Hardware.CpuMhz))

			// Look up the cluster_name for the given host_id
			clusterName := HostMapping[hostID]
			if clusterName == "" {
				clusterName = "none"
			}

			labels := prometheus.Labels{
				"host_name":    host.Name,
				"host_id":      hostID,
				"datacenter":   dc.Name,
				"cluster_name": clusterName,
			}

			cpuTotal := int64(host.Summary.Hardware.CpuMhz) * int64(host.Summary.Hardware.NumCpuCores) * int64(host.Summary.Hardware.NumCpuThreads)

			hostAvailPMem.With(labels).Set(
				float64(host.Summary.QuickStats.AvailablePMemCapacity),
			)
			hostCpuAllocRes.With(labels).Set(
				float64(*host.Config.SystemResources.Config.CpuAllocation.Reservation),
			)
			hostCpuAllocLim.With(labels).Set(
				float64(*host.Config.SystemResources.Config.CpuAllocation.Limit),
			)
			hostCpuAllocOver.With(labels).Set(
				float64(*host.Config.SystemResources.Config.CpuAllocation.OverheadLimit),
			)
			hostCpuCores.With(labels).Set(
				float64(host.Summary.Hardware.NumCpuCores),
			)
			hostCpuFree.With(labels).Set(
				float64(int64(cpuTotal) - int64(host.Summary.QuickStats.OverallCpuUsage)),
			)
			hostCpuMhz.With(labels).Set(
				float64(host.Summary.Hardware.CpuMhz),
			)
			hostCpuOverallUsage.With(labels).Set(
				float64(host.Summary.QuickStats.OverallCpuUsage),
			)
			hostCpuTotal.With(labels).Set(
				float64(cpuTotal),
			)
			hostCpuThreads.With(labels).Set(
				float64(host.Summary.Hardware.NumCpuThreads),
			)
			hostMemoryOverallUsage.With(labels).Set(
				float64(host.Summary.QuickStats.OverallMemoryUsage),
			)
			hostMemoryFree.With(labels).Set(
				float64(int64(host.Summary.Hardware.MemorySize) - (int64(host.Summary.QuickStats.OverallMemoryUsage))),
			)
			hostMemorySize.With(labels).Set(
				float64(host.Summary.Hardware.MemorySize),
			)
			hostMemoryAllocLim.With(labels).Set(
				float64(*host.Config.SystemResources.Config.MemoryAllocation.Limit),
			)
			hostMemoryAllocRes.With(labels).Set(
				float64(*host.Config.SystemResources.Config.MemoryAllocation.Reservation),
			)
			hostNicsNum.With(labels).Set(
				float64(host.Summary.Hardware.NumNics),
			)
			hostUptime.With(labels).Set(
				float64(host.Summary.QuickStats.Uptime),
			)
		}
	}
	return nil
}
