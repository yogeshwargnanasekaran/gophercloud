// +build acceptance containerinfra

package v1

import (
	"testing"

	"github.com/yogeshwargnanasekaran/gophercloud/acceptance/clients"
	"github.com/yogeshwargnanasekaran/gophercloud/acceptance/tools"
	th "github.com/yogeshwargnanasekaran/gophercloud/testhelper"
)

func TestQuotasCRUD(t *testing.T) {
	client, err := clients.NewContainerInfraV1Client()
	th.AssertNoErr(t, err)

	quota, err := CreateQuota(t, client)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, quota)
}
