// +build acceptance networking loadbalancer l7policies

package lbaas_v2

import (
	"testing"

	"github.com/yogeshwargnanasekaran/gophercloud/acceptance/clients"
	"github.com/yogeshwargnanasekaran/gophercloud/acceptance/tools"
	"github.com/yogeshwargnanasekaran/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
)

func TestL7PoliciesList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := l7policies.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list l7policies: %v", err)
	}

	allL7Policies, err := l7policies.ExtractL7Policies(allPages)
	if err != nil {
		t.Fatalf("Unable to extract l7policies: %v", err)
	}

	for _, policy := range allL7Policies {
		tools.PrintResource(t, policy)
	}
}
