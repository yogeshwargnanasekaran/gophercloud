package tokens

import "github.com/yogeshwargnanasekaran/gophercloud"

func tokenURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
