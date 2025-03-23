package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"vmware-exporter/collector"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

var (
	hostname        = os.Getenv("VSPHERE_HOSTNAME")
	insecure        = os.Getenv("VSPHERE_INSECURE") == "true"
	password        = os.Getenv("VSPHERE_PASSWORD")
	username        = os.Getenv("VSPHERE_USERNAME")
	metricsPort     = 8080
	pollingInterval = 5 * time.Minute // in minutes
)

func collectMetrics(ctx context.Context, client *govmomi.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			start := time.Now()
			err := collector.ExportClusterMetrics(ctx, client)
			if err != nil {
				log.Printf("Error exporting metrics: %v", err)
			}
			elapsed := time.Since(start)
			log.Printf("Cluster metrics retrieval took %s", elapsed)
			err = collector.ExportDatastoresMetrics(ctx, client)
			if err != nil {
				log.Printf("Error exporting metrics: %v", err)
			}
			elapsed = time.Since(start)
			log.Printf("Datastores metrics retrieval took %s", elapsed)
			start = time.Now()
			err = collector.ExportHostMetrics(ctx, client)
			if err != nil {
				log.Printf("Error exporting metrics: %v", err)
			}
			elapsed = time.Since(start)
			log.Printf("Host metrics retrieval took %s", elapsed)
			start = time.Now()
			err = collector.ExportVirtualMachineMetrics(ctx, client)
			if err != nil {
				log.Printf("Error exporting metrics: %v", err)
			}
			elapsed = time.Since(start)
			log.Printf("VM metrics retrieval took %s", elapsed)
			log.Printf("collected metrics")
			time.Sleep(pollingInterval) // Adjust the polling interval as needed
		}
	}
}

func main() {

	// Create a URL object
	u, err := soap.ParseURL(fmt.Sprintf("https://%s/sdk", hostname))
	if err != nil {
		fmt.Println("Error parsing vSphere URL:", err)
		return
	}
	u.User = url.UserPassword(username, password)

	// Create a vSphere client
	ctx := context.Background()
	client, err := govmomi.NewClient(ctx, u, insecure)
	if err != nil {
		fmt.Println("Error creating vSphere client:", err)
		return
	}
	defer client.Logout(ctx)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil)

	collectMetrics(ctx, client)

}
