package clusterconfig

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"

	"github.com/openshift/insights-operator/pkg/record"
)

// GatherOpenstackDataplanes Collects `openstackdataplanes.core.openstack.org`
// resources from all namespaces
//
// ### API Reference
// None
//
// ### Sample data
// - docs/insights-archive-sample/customresources/dataplane.openstack.org/openstackdataplanenodesets/openstack/openstack-edpm.json
//
// ### Location in archive
// - `customresources/dataplane.openstack.org/openstackdataplanes/{namespace}/{name}.json`
//
// ### Config ID
// `clusterconfig/openstack_dataplanes`
//
// ### Released version
// - 4.13
//
// ### Changes
// None
func (g *Gatherer) GatherOpenstackDataplanes(ctx context.Context) ([]record.Record, []error) {
	gatherDynamicClient, err := dynamic.NewForConfig(g.gatherKubeConfig)
	if err != nil {
		return nil, []error{err}
	}

	return gatherOpenstackDataplanes(ctx, gatherDynamicClient)
}

func gatherOpenstackDataplanes(ctx context.Context, dynamicClient dynamic.Interface) ([]record.Record, []error) {
	openstackdataplanesList, err := dynamicClient.Resource(osdpnsGroupVersionResource).List(ctx, metav1.ListOptions{})
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	var records []record.Record

	for i, osdpns := range openstackdataplanesList.Items {
		records = append(records, record.Record{
			Name: fmt.Sprintf("customresources/%s/%s/%s/%s",
				osdpnsGroupVersionResource.Group,
				osdpnsGroupVersionResource.Resource,
				osdpns.GetNamespace(),
				osdpns.GetName(),
			),
			Item: record.ResourceMarshaller{Resource: &openstackdataplanesList.Items[i]},
		})
	}

	return records, nil
}
