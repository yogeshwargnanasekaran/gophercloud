package usage

import "github.com/yogeshwargnanasekaran/gophercloud"

const resourcePath = "os-simple-tenant-usage"

func allTenantsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

func getTenantURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL(resourcePath, tenantID)
}
