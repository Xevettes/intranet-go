package grpcserver

import (
	"netbox-gateway/internal/netbox_client"

	"netbox-gateway/proto/dcim_proto"
	"netbox-gateway/proto/ipam"
	"netbox-gateway/proto/organization"
	"netbox-gateway/proto/virtualization"
)

type Server struct {
	organization.UnimplementedOrganizationServiceServer
	dcim_proto.UnimplementedDcimServiceServer
	ipam.UnimplementedIpamServiceServer
	virtualization.UnimplementedVirtualizationServiceServer
	netboxClient *netbox_client.Client
}

func NewServer(client *netbox_client.Client) *Server {
	return &Server{netboxClient: client}
}
