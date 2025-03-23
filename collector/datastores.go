package collector

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

var (
	datastoresLabels []string = []string{"datastore_name", "datastore_type"}
	dsCapacity                = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "ds",
			Name:      "capacity_bytes",
			Help:      "Datastore capacity in bytes",
		},
		datastoresLabels,
	)
	dsFreeSpace = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vmware",
			Subsystem: "ds",
			Name:      "free_bytes",
			Help:      "Datastore free space in bytes",
		},
		datastoresLabels,
	)
)

func ExportDatastoresMetrics(ctx context.Context, c *govmomi.Client) error {
	// Create a view of Datastore objects
	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datastore"}, true)
	if err != nil {
		return err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all datastores
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.Datastore.html
	var dss []mo.Datastore
	err = v.Retrieve(ctx, []string{"Datastore"}, []string{"summary"}, &dss)
	if err != nil {
		return err
	}

	for _, ds := range dss {

		populateDatastoreMapping(ds.Self.Value, ds.Summary.Name)

		labels := prometheus.Labels{
			"datastore_name": ds.Summary.Name,
			"datastore_type": ds.Summary.Type,
		}

		dsCapacity.With(labels).Set(
			float64(ds.Summary.Capacity),
		)
		dsFreeSpace.With(labels).Set(
			float64(ds.Summary.FreeSpace),
		)
	}
	return nil
}
