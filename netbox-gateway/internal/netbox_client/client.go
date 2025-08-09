package netbox_client

import (
	"github.com/netbox-community/go-netbox/v4"
)

type Client struct {
	*netbox.APIClient
}

func NewClient(netboxURL, apiToken string) (*Client, error) {
	nbClient := netbox.NewAPIClientFor(netboxURL, apiToken)
	return &Client{
		APIClient: nbClient,
	}, nil
}
