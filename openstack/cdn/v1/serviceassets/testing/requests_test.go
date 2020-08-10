package testing

import (
	"testing"

	"github.com/yogeshwargnanasekaran/gophercloud/openstack/cdn/v1/serviceassets"
	th "github.com/yogeshwargnanasekaran/gophercloud/testhelper"
	fake "github.com/yogeshwargnanasekaran/gophercloud/testhelper/client"
)

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleDeleteCDNAssetSuccessfully(t)

	err := serviceassets.Delete(fake.ServiceClient(), "96737ae3-cfc1-4c72-be88-5d0e7cc9a3f0", nil).ExtractErr()
	th.AssertNoErr(t, err)
}
