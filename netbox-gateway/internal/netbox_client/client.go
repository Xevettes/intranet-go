package netbox_client

import (
	"fmt"

	"github.com/netbox-community/go-netbox/v4"
)

type Client struct {
	*netbox.APIClient
}

func NewClient(netboxURL, apiToken string) (*Client, error) {
	if netboxURL == "" || apiToken == "" {
		return nil, fmt.Errorf("NetBox URL and API token cannot be empty")
	}
	nbClient := netbox.NewAPIClientFor(netboxURL, apiToken)
	return &Client{APIClient: nbClient}, nil
}
