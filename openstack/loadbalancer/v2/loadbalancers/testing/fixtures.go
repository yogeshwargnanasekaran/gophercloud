package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/yogeshwargnanasekaran/gophercloud/openstack/loadbalancer/v2/listeners"
	"github.com/yogeshwargnanasekaran/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/yogeshwargnanasekaran/gophercloud/openstack/loadbalancer/v2/monitors"
	"github.com/yogeshwargnanasekaran/gophercloud/openstack/loadbalancer/v2/pools"
	th "github.com/yogeshwargnanasekaran/gophercloud/testhelper"
	"github.com/yogeshwargnanasekaran/gophercloud/testhelper/client"
)

// LoadbalancersListBody contains the canned body of a loadbalancer list response.
const LoadbalancersListBody = `
{
	"loadbalancers":[
	         {
			"id": "c331058c-6a40-4144-948e-b9fb1df9db4b",
			"project_id": "54030507-44f7-473c-9342-b4d14a95f692",
			"created_at": "2019-06-30T04:15:37",
			"updated_at": "2019-06-30T05:18:49",
			"name": "web_lb",
			"description": "lb config for the web tier",
			"vip_subnet_id": "8a49c438-848f-467b-9655-ea1548708154",
			"vip_address": "10.30.176.47",
			"vip_port_id": "2a22e552-a347-44fd-b530-1f2b1b2a6735",
			"flavor_id": "60df399a-ee85-11e9-81b4-2a2ae2dbcce4",
			"provider": "haproxy",
			"admin_state_up": true,
			"provisioning_status": "ACTIVE",
			"operating_status": "ONLINE",
			"tags": ["test", "stage"]
		},
		{
			"id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
			"project_id": "54030507-44f7-473c-9342-b4d14a95f692",
			"created_at": "2019-06-30T04:15:37",
			"updated_at": "2019-06-30T05:18:49",
			"name": "db_lb",
			"description": "lb config for the db tier",
			"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
			"vip_address": "10.30.176.48",
			"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
			"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
			"provider": "haproxy",
			"admin_state_up": true,
			"provisioning_status": "PENDING_CREATE",
			"operating_status": "OFFLINE",
			"tags": ["test", "stage"]
		}
	]
}
`

// SingleLoadbalancerBody is the canned body of a Get request on an existing loadbalancer.
const SingleLoadbalancerBody = `
{
	"loadbalancer": {
		"id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
		"project_id": "54030507-44f7-473c-9342-b4d14a95f692",
		"created_at": "2019-06-30T04:15:37",
		"updated_at": "2019-06-30T05:18:49",
		"name": "db_lb",
		"description": "lb config for the db tier",
		"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		"vip_address": "10.30.176.48",
		"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		"provider": "haproxy",
		"admin_state_up": true,
		"provisioning_status": "PENDING_CREATE",
		"operating_status": "OFFLINE",
		"tags": ["test", "stage"]
	}
}
`

// PostUpdateLoadbalancerBody is the canned response body of a Update request on an existing loadbalancer.
const PostUpdateLoadbalancerBody = `
{
	"loadbalancer": {
		"id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
		"project_id": "54030507-44f7-473c-9342-b4d14a95f692",
		"created_at": "2019-06-30T04:15:37",
		"updated_at": "2019-06-30T05:18:49",
		"name": "NewLoadbalancerName",
		"description": "lb config for the db tier",
		"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		"vip_address": "10.30.176.48",
		"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		"provider": "haproxy",
		"admin_state_up": true,
		"provisioning_status": "PENDING_CREATE",
		"operating_status": "OFFLINE",
		"tags": ["test"]
	}
}
`

// GetLoadbalancerStatusesBody is the canned request body of a Get request on loadbalancer's status.
const GetLoadbalancerStatusesBody = `
{
	"statuses" : {
		"loadbalancer": {
			"id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
			"name": "db_lb",
			"provisioning_status": "PENDING_UPDATE",
			"operating_status": "ACTIVE",
			"tags": ["test", "stage"],
			"listeners": [{
				"id": "db902c0c-d5ff-4753-b465-668ad9656918",
				"name": "db",
				"provisioning_status": "ACTIVE",
				"pools": [{
					"id": "fad389a3-9a4a-4762-a365-8c7038508b5d",
					"name": "db",
					"provisioning_status": "ACTIVE",
					"healthmonitor": {
						"id": "67306cda-815d-4354-9fe4-59e09da9c3c5",
						"type":"PING",
						"provisioning_status": "ACTIVE"
					},
					"members":[{
						"id": "2a280670-c202-4b0b-a562-34077415aabf",
						"name": "db",
						"address": "10.0.2.11",
						"protocol_port": 80,
						"provisioning_status": "ACTIVE"
					}]
				}]
			}]
		}
	}
}
`

// LoadbalancerStatsTree is the canned request body of a Get request on loadbalancer's statistics.
const GetLoadbalancerStatsBody = `
{
    "stats": {
        "active_connections": 0,
        "bytes_in": 9532,
        "bytes_out": 22033,
        "request_errors": 46,
        "total_connections": 112
    }
}
`

var createdTime, _ = time.Parse(time.RFC3339, "2019-06-30T04:15:37Z")
var updatedTime, _ = time.Parse(time.RFC3339, "2019-06-30T05:18:49Z")

var (
	LoadbalancerWeb = loadbalancers.LoadBalancer{
		ID:                 "c331058c-6a40-4144-948e-b9fb1df9db4b",
		ProjectID:          "54030507-44f7-473c-9342-b4d14a95f692",
		CreatedAt:          createdTime,
		UpdatedAt:          updatedTime,
		Name:               "web_lb",
		Description:        "lb config for the web tier",
		VipSubnetID:        "8a49c438-848f-467b-9655-ea1548708154",
		VipAddress:         "10.30.176.47",
		VipPortID:          "2a22e552-a347-44fd-b530-1f2b1b2a6735",
		FlavorID:           "60df399a-ee85-11e9-81b4-2a2ae2dbcce4",
		Provider:           "haproxy",
		AdminStateUp:       true,
		ProvisioningStatus: "ACTIVE",
		OperatingStatus:    "ONLINE",
		Tags:               []string{"test", "stage"},
	}
	LoadbalancerDb = loadbalancers.LoadBalancer{
		ID:                 "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
		ProjectID:          "54030507-44f7-473c-9342-b4d14a95f692",
		CreatedAt:          createdTime,
		UpdatedAt:          updatedTime,
		Name:               "db_lb",
		Description:        "lb config for the db tier",
		VipSubnetID:        "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		VipAddress:         "10.30.176.48",
		VipPortID:          "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		FlavorID:           "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		Provider:           "haproxy",
		AdminStateUp:       true,
		ProvisioningStatus: "PENDING_CREATE",
		OperatingStatus:    "OFFLINE",
		Tags:               []string{"test", "stage"},
	}
	LoadbalancerUpdated = loadbalancers.LoadBalancer{
		ID:                 "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
		ProjectID:          "54030507-44f7-473c-9342-b4d14a95f692",
		CreatedAt:          createdTime,
		UpdatedAt:          updatedTime,
		Name:               "NewLoadbalancerName",
		Description:        "lb config for the db tier",
		VipSubnetID:        "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		VipAddress:         "10.30.176.48",
		VipPortID:          "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		FlavorID:           "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		Provider:           "haproxy",
		AdminStateUp:       true,
		ProvisioningStatus: "PENDING_CREATE",
		OperatingStatus:    "OFFLINE",
		Tags:               []string{"test"},
	}
	LoadbalancerStatusesTree = loadbalancers.StatusTree{
		Loadbalancer: &loadbalancers.LoadBalancer{
			ID:                 "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
			Name:               "db_lb",
			ProvisioningStatus: "PENDING_UPDATE",
			OperatingStatus:    "ACTIVE",
			Tags:               []string{"test", "stage"},
			Listeners: []listeners.Listener{{
				ID:                 "db902c0c-d5ff-4753-b465-668ad9656918",
				Name:               "db",
				ProvisioningStatus: "ACTIVE",
				Pools: []pools.Pool{{
					ID:                 "fad389a3-9a4a-4762-a365-8c7038508b5d",
					Name:               "db",
					ProvisioningStatus: "ACTIVE",
					Monitor: monitors.Monitor{
						ID:                 "67306cda-815d-4354-9fe4-59e09da9c3c5",
						Type:               "PING",
						ProvisioningStatus: "ACTIVE",
					},
					Members: []pools.Member{{
						ID:                 "2a280670-c202-4b0b-a562-34077415aabf",
						Name:               "db",
						Address:            "10.0.2.11",
						ProtocolPort:       80,
						ProvisioningStatus: "ACTIVE",
					}},
				}},
			}},
		},
	}
	LoadbalancerStatsTree = loadbalancers.Stats{
		ActiveConnections: 0,
		BytesIn:           9532,
		BytesOut:          22033,
		RequestErrors:     46,
		TotalConnections:  112,
	}
)

// HandleLoadbalancerListSuccessfully sets up the test server to respond to a loadbalancer List request.
func HandleLoadbalancerListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, LoadbalancersListBody)
		case "45e08a3e-a78f-4b40-a229-1e7e23eee1ab":
			fmt.Fprintf(w, `{ "loadbalancers": [] }`)
		default:
			t.Fatalf("/v2.0/lbaas/loadbalancers invoked with unexpected marker=[%s]", marker)
		}
	})
}

// HandleLoadbalancerCreationSuccessfully sets up the test server to respond to a loadbalancer creation request
// with a given response.
func HandleLoadbalancerCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"loadbalancer": {
				"name": "db_lb",
				"vip_port_id": "2bf413c8-41a9-4477-b505-333d5cbe8b55",
				"vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
				"vip_address": "10.30.176.48",
				"flavor_id": "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
				"provider": "haproxy",
				"admin_state_up": true,
				"tags": ["test", "stage"]
			}
		}`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

// HandleLoadbalancerGetSuccessfully sets up the test server to respond to a loadbalancer Get request.
func HandleLoadbalancerGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, SingleLoadbalancerBody)
	})
}

// HandleLoadbalancerGetStatusesTree sets up the test server to respond to a loadbalancer Get statuses tree request.
func HandleLoadbalancerGetStatusesTree(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab/status", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, GetLoadbalancerStatusesBody)
	})
}

// HandleLoadbalancerDeletionSuccessfully sets up the test server to respond to a loadbalancer deletion request.
func HandleLoadbalancerDeletionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleLoadbalancerUpdateSuccessfully sets up the test server to respond to a loadbalancer Update request.
func HandleLoadbalancerUpdateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `{
			"loadbalancer": {
				"name": "NewLoadbalancerName",
				"tags": ["test"]
			}
		}`)

		fmt.Fprintf(w, PostUpdateLoadbalancerBody)
	})
}

// HandleLoadbalancerGetStatsTree sets up the test server to respond to a loadbalancer Get stats tree request.
func HandleLoadbalancerGetStatsTree(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab/stats", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, GetLoadbalancerStatsBody)
	})
}

// HandleLoadbalancerFailoverSuccessfully sets up the test server to respond to a loadbalancer failover request.
func HandleLoadbalancerFailoverSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/36e08a3e-a78f-4b40-a229-1e7e23eee1ab/failover", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})
}
