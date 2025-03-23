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
	vmLabels          []string = []string{"machine_name", "datacenter", "cluster_name"}
	vmDatastoreLabels []string = []string{"machine_name", "host_name", "datacenter", "cluster_name", "datastore_id", "datastore_name"}
	vmDiskLabels      []string = []string{"machine_name", "host_name", "datacenter", "cluster_name", "disk_path"}
	vmCpuAllocLim              = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_allocation_limit_mhz",
			Help:      "VM CPU allocation limit in Mhz",
		},
		vmLabels,
	)
	vmCpuAllocRes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_allocation_reservation_mhz",
			Help:      "VM CPU allocation reservation in Mhz",
		},
		vmLabels,
	)
	vmCpuEnt = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_entitled_bytes",
			Help:      "VM entitled cpu in mhz",
		},
		vmLabels,
	)
	vmCpuMhz = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_mhz",
			Help:      "VM CPU core Mhz",
		},
		vmLabels,
	)
	vmCpuReservation = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_reservation_mhz",
			Help:      "VM CPU reservation in Mhz",
		},
		vmLabels,
	)
	vmCpuOverallDemand = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_usage_mhz",
			Help:      "Overall VM CPU demand in Mhz",
		},
		vmLabels,
	)

	vmCpuMaxUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_usage_max",
			Help:      "Max VM CPU usage in Mhz",
		},
		vmLabels,
	)
	vmCpuNum = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_cores_total",
			Help:      "VM CPU number of cores",
		},
		vmLabels,
	)
	vmCpuOverallUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "cpu_mhz_total",
			Help:      "Overall VM CPU usage in Mhz",
		},
		vmLabels,
	)
	vmCreationDate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "creation_date_seconds",
			Help:      "VM creation date in seconds",
		},
		vmLabels,
	)
	vmDatastoreCommited = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "datastore_committed_bytes",
			Help:      "VM committed storage in bytes",
		},
		vmDatastoreLabels,
	)
	vmDatastoreUncommited = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "datastore_uncommitted_bytes",
			Help:      "VM uncommitted storage in bytes",
		},
		vmDatastoreLabels,
	)
	vmDiskCapacity = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "disk_capacity_bytes",
			Help:      "VM disk capacity in bytes",
		},
		vmDiskLabels,
	)
	vmDiskFreeSpace = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "disk_free_space_bytes",
			Help:      "VM disk free space in bytes",
		},
		vmDiskLabels,
	)
	vmDiskMappingKey = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "disk_mapping_key",
			Help:      "VM disk mapping key",
		},
		vmDiskLabels,
	)
	vmMemoryActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "memory_active_bytes",
			Help:      "VM active memory in bytes",
		},
		vmLabels,
	)
	vmMemoryAllocLim = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "memory_allocation_limit_bytes",
			Help:      "VM memory allocation limit in bytes",
		},
		vmLabels,
	)
	vmMemoryAllocRes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "memory_allocation_reservation_bytes",
			Help:      "VM memory allocation reservation in bytes",
		},
		vmLabels,
	)
	vmMemoryGranted = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "memory_granted_bytes",
			Help:      "VM granted memory in bytes",
		},
		vmLabels,
	)
	vmMemoryReservation = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "memory_reservation_bytes",
			Help:      "VM memory reservation in bytes",
		},
		vmLabels,
	)
	vmMemoryUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "memory_used_bytes",
			Help:      "VM used memory in bytes",
		},
		vmLabels,
	)
	vmMemoryEnt = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "memory_entitled_bytes",
			Help:      "VM entitled memory in bytes",
		},
		vmLabels,
	)
	vmMemoryTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "memory_bytes_total",
			Help:      "VM total memory in bytes",
		},
		vmLabels,
	)
	vmStorageCommited = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "storage_committed_bytes",
			Help:      "VM storage committed in bytes",
		},
		vmLabels,
	)
	vmUptime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "vm",
			Name:      "uptime_seconds",
			Help:      "VM uptime in seconds",
		},
		vmLabels,
	)
)

func ExportVirtualMachineMetrics(ctx context.Context, client *govmomi.Client) error {

	c := client.Client

	// Create a view manager
	m := view.NewManager(client.Client)

	// Create a container view for the datacenters
	containerView, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		fmt.Printf("Error creating container view: %v\n", err)
		return err
	}
	defer containerView.Destroy(ctx)

	// Retrieve a list of datacenters
	var datacenters []mo.Datacenter
	err = containerView.Retrieve(ctx, []string{"Datacenter"}, nil, &datacenters)
	if err != nil {
		fmt.Printf("Error retrieving datacenters: %v\n", err)
		return err
	}

	// Iterate through the datacenters and list vms
	for _, dc := range datacenters {

		// Create a container view for the vms in the datacenter
		containerView, err := m.CreateContainerView(ctx, dc.Reference(), []string{"VirtualMachine"}, true)
		if err != nil {
			fmt.Printf("Error creating container view for vms: %v\n", err)
			return err
		}
		defer containerView.Destroy(ctx)

		// Retrieve a list of vms in the datacenter
		var vms []mo.VirtualMachine
		err = containerView.Retrieve(ctx, []string{"VirtualMachine"}, nil, &vms)
		if err != nil {
			fmt.Printf("Error retrieving vms: %v\n", err)
			return err
		}

		for _, vm := range vms {
			hostID := vm.Summary.Runtime.Host.Value

			clusterName := HostMapping[hostID]
			if clusterName == "" {
				clusterName = "none"
			}

			hostName := VirtualMachineMapping[hostID]
			if hostName == "" {
				hostName = "unknown"
			}

			vmCpuSpeed := HostConfig[hostName]

			labels := prometheus.Labels{
				"machine_name": vm.Name,
				"datacenter":   dc.Name,
				"cluster_name": clusterName,
			}

			// collect datastore metrics
			for _, storage := range vm.Storage.PerDatastoreUsage {
				datastoreId := storage.Datastore.Value
				datastoreCommitted := storage.Committed
				datastoreUncommitted := storage.Uncommitted

				datastoreName := DatastoreMapping[datastoreId]
				if datastoreName == "" {
					datastoreName = "unknown"
				}

				datastoresLabels := prometheus.Labels{
					"machine_name":   vm.Name,
					"host_name":      hostName,
					"datacenter":     dc.Name,
					"cluster_name":   clusterName,
					"datastore_id":   datastoreId,
					"datastore_name": datastoreName,
				}
				vmDatastoreCommited.With(datastoresLabels).Set(
					float64(datastoreCommitted),
				)
				vmDatastoreUncommited.With(datastoresLabels).Set(
					float64(datastoreUncommitted),
				)
			}

			for _, disk := range vm.Guest.Disk {
				vmDiskLabels := prometheus.Labels{
					"machine_name": vm.Name,
					"host_name":    hostName,
					"datacenter":   dc.Name,
					"cluster_name": clusterName,
					"disk_path":    disk.DiskPath,
				}
				for _, mapping := range disk.Mappings {
					vmDiskMappingKey.With(vmDiskLabels).Set(
						float64(mapping.Key),
					)
				}
				vmDiskCapacity.With(vmDiskLabels).Set(
					float64(disk.Capacity),
				)
				vmDiskFreeSpace.With(vmDiskLabels).Set(
					float64(disk.FreeSpace),
				)
			}

			vmCpuAllocLim.With(labels).Set(
				float64(*vm.Config.CpuAllocation.Limit),
			)
			vmCpuAllocRes.With(labels).Set(
				float64(*vm.Config.CpuAllocation.Reservation),
			)
			vmCpuEnt.With(labels).Set(
				float64(vm.Summary.QuickStats.StaticCpuEntitlement),
			)
			vmCpuMaxUsage.With(labels).Set(
				float64(vm.Summary.Runtime.MaxCpuUsage),
			)
			vmCpuMhz.With(labels).Set(
				float64(vmCpuSpeed),
			)
			vmCpuNum.With(labels).Set(
				float64(vm.Config.Hardware.NumCPU),
			)
			vmCpuOverallDemand.With(labels).Set(
				float64(vm.Summary.QuickStats.OverallCpuDemand),
			)
			vmCpuOverallUsage.With(labels).Set(
				float64(vm.Summary.QuickStats.OverallCpuUsage),
			)
			vmCpuReservation.With(labels).Set(
				float64(vm.Summary.Config.CpuReservation),
			)
			vmCreationDate.With(labels).Set(
				float64(vm.Config.CreateDate.Unix()),
			)
			vmMemoryActive.With(labels).Set(
				float64(vm.Summary.QuickStats.ActiveMemory),
			)
			vmMemoryAllocLim.With(labels).Set(
				float64(*vm.Config.MemoryAllocation.Limit),
			)
			vmMemoryAllocRes.With(labels).Set(
				float64(*vm.Config.MemoryAllocation.Reservation),
			)
			vmMemoryEnt.With(labels).Set(
				float64(vm.Summary.QuickStats.StaticMemoryEntitlement),
			)
			vmMemoryGranted.With(labels).Set(
				float64(vm.Summary.QuickStats.GrantedMemory),
			)
			vmMemoryReservation.With(labels).Set(
				float64(vm.Summary.Config.MemoryReservation),
			)
			vmMemoryTotal.With(labels).Set(
				float64(vm.Config.Hardware.MemoryMB),
			)
			vmMemoryUsage.With(labels).Set(
				float64(vm.Summary.QuickStats.GuestMemoryUsage),
			)
			vmStorageCommited.With(labels).Set(
				float64(vm.Summary.Storage.Committed),
			)
			vmUptime.With(labels).Set(
				float64(vm.Summary.QuickStats.UptimeSeconds),
			)
		}
	}
	return nil
}
