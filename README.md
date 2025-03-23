# VMware Exporter

VMware Exporter is a Go-based application that collects metrics from a VMware vSphere environment and exposes them in Prometheus format. It gathers metrics for clusters, hosts, datastores, and virtual machines, providing insights into resource usage and performance.

## Features

- Collects metrics for:
  - Clusters (CPU, memory, hosts, threads, etc.)
  - Hosts (CPU, memory, NICs, uptime, etc.)
  - Datastores (capacity, free space, etc.)
  - Virtual Machines (CPU, memory, storage, uptime, etc.)
- Exposes metrics in Prometheus format via an HTTP endpoint.
- Supports configurable polling intervals.

## Prerequisites

- VMware vSphere environment.
- Prometheus server for scraping metrics.
- Docker (optional, for containerized deployment).

## Installation

### Clone the Repository

```bash
git clone https://github.com/your-repo/vmware-exporter.git
cd vmware-exporter
```

### Run with Docker

Build and run the Docker container:

```bash
docker build -t vmware-exporter .
docker run -e VSPHERE_HOSTNAME=<hostname> -e VSPHERE_USERNAME=<username> -e VSPHERE_PASSWORD=<password> -e VSPHERE_INSECURE=true -p 8080:8080 vmware-exporter
```

## Configuration

The application uses the following environment variables for configuration:

- `VSPHERE_HOSTNAME`: The hostname or IP address of the vSphere server.
- `VSPHERE_USERNAME`: The username for vSphere authentication.
- `VSPHERE_PASSWORD`: The password for vSphere authentication.
- `VSPHERE_INSECURE`: Set to `true` to allow insecure connections (default: `false`).
- `METRICS_PORT`: The port to expose metrics (default: `8080`).
- `POLLING_INTERVAL`: The interval for polling metrics (default: `5m`).

## Usage

Run the application:

```bash
export VSPHERE_HOSTNAME=<hostname>
export VSPHERE_USERNAME=<username>
export VSPHERE_PASSWORD=<password>
export VSPHERE_INSECURE=true
./vmware-exporter
```

Access the metrics at `http://<your-server>:8080/metrics`.

## Metrics

These metrics are exposed in Prometheus format via the `/metrics` HTTP endpoint.

### Cluster Metrics
- `vmware_cluster_cpu_effective_mhz`: Effective Cluster CPU in MHz.
- `vmware_cluster_cpu_cores_total`: Total Cluster CPU cores.
- `vmware_cluster_cpu_mhz_total`: Total Cluster CPU in MHz.
- `vmware_cluster_hosts_effective_total`: Effective Cluster hosts.
- `vmware_cluster_hosts_total`: Total Cluster hosts.
- `vmware_cluster_memory_effective_bytes`: Effective Cluster memory in bytes.
- `vmware_cluster_memory_bytes_total`: Total Cluster memory in bytes.
- `vmware_cluster_threads_total`: Total Cluster threads.

### Host Metrics
- `vmware_host_available_pmem_bytes`: Host available persistent memory in bytes.
- `vmware_host_cpu_allocation_reservation_mhz`: Host CPU allocation reservation in MHz.
- `vmware_host_cpu_allocation_limit_mhz`: Host CPU allocation limit in MHz.
- `vmware_host_cpu_allocation_overhead_mhz`: Host CPU allocation overhead in MHz.
- `vmware_host_cpu_cores_total`: Total Host CPU cores.
- `vmware_host_cpu_free_mhz`: Free Host CPU in MHz.
- `vmware_host_cpu_usage_mhz`: Overall Host CPU usage in MHz.
- `vmware_host_cpu_core_mhz`: Host CPU core MHz.
- `vmware_host_cpu_mhz_total`: Total Host CPU in MHz.
- `vmware_host_cpu_threads_total`: Total Host CPU threads.
- `vmware_host_memory_allocation_bytes`: Host memory allocation in bytes.
- `vmware_host_memory_allocation_limit_bytes`: Host memory allocation limit in bytes.
- `vmware_host_memory_free_bytes`: Host memory free in bytes.
- `vmware_host_memory_bytes_total`: Total Host memory in bytes.
- `vmware_host_memory_usage_bytes`: Overall Host memory usage in bytes.
- `vmware_host_nics_total`: Total Host NICs.
- `vmware_host_uptime_seconds`: Host uptime in seconds.

### Datastore Metrics
- `vmware_ds_capacity_bytes`: Datastore capacity in bytes.
- `vmware_ds_free_bytes`: Datastore free space in bytes.

### Virtual Machine Metrics
- `vmware_vm_cpu_allocation_limit_mhz`: VM CPU allocation limit in MHz.
- `vmware_vm_cpu_allocation_reservation_mhz`: VM CPU allocation reservation in MHz.
- `vmware_vm_cpu_entitled_bytes`: VM entitled CPU in MHz.
- `vmware_vm_cpu_mhz`: VM CPU core MHz.
- `vmware_vm_cpu_usage_mhz`: Overall VM CPU demand in MHz.
- `vmware_vm_cpu_usage_max`: Max VM CPU usage in MHz.
- `vmware_vm_cpu_cores_total`: VM CPU number of cores.
- `vmware_vm_memory_active_bytes`: VM active memory in bytes.
- `vmware_vm_memory_allocation_limit_bytes`: VM memory allocation limit in bytes.
- `vmware_vm_memory_allocation_reservation_bytes`: VM memory allocation reservation in bytes.
- `vmware_vm_memory_granted_bytes`: VM granted memory in bytes.
- `vmware_vm_memory_reservation_bytes`: VM memory reservation in bytes.
- `vmware_vm_memory_bytes_total`: VM total memory in bytes.
- `vmware_vm_memory_usage_bytes`: VM used memory in bytes.
- `vmware_vm_memory_entitled_bytes`: VM entitled memory in bytes.
- `vmware_vm_storage_committed_bytes`: VM storage committed in bytes.
- `vmware_vm_uptime_seconds`: VM uptime in seconds.
- `vmware_vm_creation_date_seconds`: VM creation date in seconds.
- `vmware_vm_datastore_committed_bytes`: VM committed storage in bytes.
- `vmware_vm_datastore_uncommitted_bytes`: VM uncommitted storage in bytes.
- `vmware_vm_disk_capacity_bytes`: VM disk capacity in bytes.
- `vmware_vm_disk_free_space_bytes`: VM disk free space in bytes.
- `vmware_vm_disk_mapping_key`: VM disk mapping key.

### Dependencies

This project uses the following dependencies:

- [govmomi](https://github.com/vmware/govmomi) for interacting with VMware vSphere.
- [Prometheus client_golang](https://github.com/prometheus/client_golang) for exposing metrics.

### Run Locally

Install dependencies:

```bash
go mod tidy
```

Run the application:

```bash
go run 

main.go


```

### Testing

Unit tests can be added to validate functionality. Use the following command to run tests:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [VMware vSphere SDK](https://developer.vmware.com/web/sdk/60/vsphere-vim)
- [Prometheus](https://prometheus.io/)
```