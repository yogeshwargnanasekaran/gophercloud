// +build acceptance

package v1

import (
	"strings"
	"testing"

	"github.com/yogeshwargnanasekaran/gophercloud/acceptance/clients"
	"github.com/yogeshwargnanasekaran/gophercloud/acceptance/tools"
	"github.com/yogeshwargnanasekaran/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/yogeshwargnanasekaran/gophercloud/pagination"
	th "github.com/yogeshwargnanasekaran/gophercloud/testhelper"
)

// numContainers is the number of containers to create for testing.
var numContainers = 2

func TestContainers(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	// Create a slice of random container names.
	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = tools.RandomString("gophercloud-test-container-", 8)
	}

	// Create numContainers containers.
	for i := 0; i < len(cNames); i++ {
		res := containers.Create(client, cNames[i], nil)
		th.AssertNoErr(t, res.Err)
	}
	// Delete the numContainers containers after function completion.
	defer func() {
		for i := 0; i < len(cNames); i++ {
			res := containers.Delete(client, cNames[i])
			th.AssertNoErr(t, res.Err)
		}
	}()

	// List the numContainer names that were just created. To just list those,
	// the 'prefix' parameter is used.
	err = containers.List(client, &containers.ListOpts{Full: true, Prefix: "gophercloud-test-container-"}).EachPage(func(page pagination.Page) (bool, error) {
		containerList, err := containers.ExtractInfo(page)
		th.AssertNoErr(t, err)

		for _, n := range containerList {
			t.Logf("Container: Name [%s] Count [%d] Bytes [%d]",
				n.Name, n.Count, n.Bytes)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)

	// List the info for the numContainer containers that were created.
	err = containers.List(client, &containers.ListOpts{Full: false, Prefix: "gophercloud-test-container-"}).EachPage(func(page pagination.Page) (bool, error) {
		containerList, err := containers.ExtractNames(page)
		th.AssertNoErr(t, err)
		for _, n := range containerList {
			t.Logf("Container: Name [%s]", n)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)

	// Update one of the numContainer container metadata.
	metadata := map[string]string{
		"Gophercloud-Test": "containers",
	}

	updateres := containers.Update(client, cNames[0], &containers.UpdateOpts{Metadata: metadata})
	th.AssertNoErr(t, updateres.Err)
	// After the tests are done, delete the metadata that was set.
	defer func() {
		temp := []string{}
		for k, _ := range metadata {
			temp = append(temp, k)
		}
		res := containers.Update(client, cNames[0], &containers.UpdateOpts{RemoveMetadata: temp})
		th.AssertNoErr(t, res.Err)

		// confirm the metadata was removed
		getOpts := containers.GetOpts{
			Newest: true,
		}

		cm, err := containers.Get(client, cNames[0], getOpts).ExtractMetadata()
		th.AssertNoErr(t, err)
		for k, _ := range metadata {
			if _, ok := cm[k]; ok {
				t.Errorf("Unexpected custom metadata with key: %s", k)
			}
		}
	}()

	// Retrieve a container's metadata.
	getOpts := containers.GetOpts{
		Newest: true,
	}

	cm, err := containers.Get(client, cNames[0], getOpts).ExtractMetadata()
	th.AssertNoErr(t, err)
	for k := range metadata {
		if cm[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
		}
	}
}

func TestListAllContainers(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	numContainers := 20

	// Create a slice of random container names.
	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = tools.RandomString("gophercloud-test-container-", 8)
	}

	// Create numContainers containers.
	for i := 0; i < len(cNames); i++ {
		res := containers.Create(client, cNames[i], nil)
		th.AssertNoErr(t, res.Err)
	}
	// Delete the numContainers containers after function completion.
	defer func() {
		for i := 0; i < len(cNames); i++ {
			res := containers.Delete(client, cNames[i])
			th.AssertNoErr(t, res.Err)
		}
	}()

	// List all the numContainer names that were just created. To just list those,
	// the 'prefix' parameter is used.
	allPages, err := containers.List(client, &containers.ListOpts{Full: true, Limit: 5, Prefix: "gophercloud-test-container-"}).AllPages()
	th.AssertNoErr(t, err)
	containerInfoList, err := containers.ExtractInfo(allPages)
	th.AssertNoErr(t, err)
	for _, n := range containerInfoList {
		t.Logf("Container: Name [%s] Count [%d] Bytes [%d]",
			n.Name, n.Count, n.Bytes)
	}
	th.AssertEquals(t, numContainers, len(containerInfoList))

	// List the info for all the numContainer containers that were created.
	allPages, err = containers.List(client, &containers.ListOpts{Full: false, Limit: 2, Prefix: "gophercloud-test-container-"}).AllPages()
	th.AssertNoErr(t, err)
	containerNamesList, err := containers.ExtractNames(allPages)
	th.AssertNoErr(t, err)
	for _, n := range containerNamesList {
		t.Logf("Container: Name [%s]", n)
	}
	th.AssertEquals(t, numContainers, len(containerNamesList))
}

func TestBulkDeleteContainers(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	numContainers := 20

	// Create a slice of random container names.
	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = tools.RandomString("test&happy?-", 8)
	}

	// Create numContainers containers.
	for i := 0; i < len(cNames); i++ {
		res := containers.Create(client, cNames[i], nil)
		th.AssertNoErr(t, res.Err)
	}

	expectedResp := containers.BulkDeleteResponse{
		ResponseStatus: "200 OK",
		Errors:         [][]string{},
		NumberDeleted:  numContainers,
	}

	resp, err := containers.BulkDelete(client, cNames).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, *resp)
	th.AssertDeepEquals(t, *resp, expectedResp)

	for _, c := range cNames {
		_, err = containers.Get(client, c, nil).Extract()
		th.AssertErr(t, err)
	}
}
