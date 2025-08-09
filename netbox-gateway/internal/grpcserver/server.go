package grpcserver

import (
	"netbox-gateway/internal/netbox_client"

	// Importando todos os pacotes proto necessários
	"netbox-gateway/proto/dcim"
	"netbox-gateway/proto/ipam"
	"netbox-gateway/proto/organization"
	"netbox-gateway/proto/virtualization"
)

// Server implementa todas as interfaces de serviço gRPC para o NetBox.
// Ao "embedar" os Unimplemented...Servers, garantimos compatibilidade futura.
type Server struct {
	organization.UnimplementedOrganizationServiceServer
	dcim.UnimplementedDcimServiceServer
	ipam.UnimplementedIpamServiceServer
	virtualization.UnimplementedVirtualizationServiceServer

	netboxClient *netbox_client.Client
}

// NewServer cria uma nova instância do nosso servidor gRPC.
func NewServer(client *netbox_client.Client) *Server {
	return &Server{netboxClient: client}
}
