package collector

var (
	HostConfig            = make(map[string]float64)
	HostMapping           = make(map[string]string)
	VirtualMachineMapping = make(map[string]string)
	DatastoreMapping      = make(map[string]string)
)

func populateDatastoreMapping(datastoreID, datastoreName string) {
	DatastoreMapping[datastoreID] = datastoreName
}

func populateHostConfig(hostName string, cpuMhz float64) {
	HostConfig[hostName] = cpuMhz
}

func populateHostMapping(hostID, clusterName string) {
	HostMapping[hostID] = clusterName
}

func populateVirtualMachineMapping(hostID, hostName string) {
	VirtualMachineMapping[hostID] = hostName
}
