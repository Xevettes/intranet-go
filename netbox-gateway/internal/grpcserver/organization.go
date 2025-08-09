package grpcserver

import (
	"context"
	"netbox-gateway/proto/organization"

	"github.com/go-openapi/swag"
	netboxClientDcim "github.com/netbox-community/go-netbox/v4/netbox/dcim"
	"github.com/netbox-community/go-netbox/v4/netbox/models"
	netboxClientTenancy "github.com/netbox-community/go-netbox/v4/netbox/tenancy"
)

// --- Implementação do Serviço: Organization ---

func (s *Server) ListTenants(ctx context.Context, req *organization.ListRequest) (*organization.ListTenantsResponse, error) {
	params := netboxClientTenancy.NewTenancyTenantsListParams().
		WithLimit(swag.Int64(req.GetLimit())).
		WithOffset(swag.Int64(req.GetOffset()))

	result, err := s.netboxClient.Tenancy.TenancyTenantsList(params, nil)
	if err != nil {
		return nil, err
	}

	protoTenants := make([]*organization.Tenant, len(result.Payload.Results))
	for i, tenant := range result.Payload.Results {
		protoTenants[i] = &organization.Tenant{
			Id:          tenant.ID,
			Name:        *tenant.Name,
			Slug:        *tenant.Slug,
			Description: tenant.Description,
		}
	}

	return &organization.ListTenantsResponse{
		Results: protoTenants,
		Total:   swag.Int64Value(result.Payload.Count),
	}, nil
}

func (s *Server) GetTenant(ctx context.Context, req *organization.GetRequest) (*organization.Tenant, error) {
	params := netboxClientTenancy.NewTenancyTenantsReadParams().WithID(req.GetId())
	result, err := s.netboxClient.Tenancy.TenancyTenantsRead(params, nil)
	if err != nil {
		return nil, err
	}
	tenant := result.Payload
	return &organization.Tenant{
		Id:          tenant.ID,
		Name:        *tenant.Name,
		Slug:        *tenant.Slug,
		Description: tenant.Description,
	}, nil
}

func (s *Server) CreateTenant(ctx context.Context, req *organization.CreateTenantRequest) (*organization.Tenant, error) {
	params := netboxClientTenancy.NewTenancyTenantsCreateParams().WithData(&models.WritableTenant{
		Name:        &req.Name,
		Slug:        &req.Slug,
		Description: req.Description,
	})

	result, err := s.netboxClient.Tenancy.TenancyTenantsCreate(params, nil)
	if err != nil {
		return nil, err
	}
	tenant := result.Payload
	return &organization.Tenant{
		Id:          tenant.ID,
		Name:        *tenant.Name,
		Slug:        *tenant.Slug,
		Description: tenant.Description,
	}, nil
}

func (s *Server) ListSites(ctx context.Context, req *organization.ListRequest) (*organization.ListSitesResponse, error) {
	params := netboxClientDcim.NewDcimSitesListParams().
		WithLimit(swag.Int64(req.GetLimit())).
		WithOffset(swag.Int64(req.GetOffset()))

	result, err := s.netboxClient.Dcim.DcimSitesList(params, nil)
	if err != nil {
		return nil, err
	}

	protoSites := make([]*organization.Site, len(result.Payload.Results))
	for i, site := range result.Payload.Results {
		protoSites[i] = &organization.Site{
			Id:     site.ID,
			Name:   *site.Name,
			Slug:   *site.Slug,
			Status: swag.StringValue(site.Status.Value),
		}
	}

	return &organization.ListSitesResponse{
		Results: protoSites,
		Total:   swag.Int64Value(result.Payload.Count),
	}, nil
}

func (s *Server) GetSite(ctx context.Context, req *organization.GetRequest) (*organization.Site, error) {
	params := netboxClientDcim.NewDcimSitesReadParams().WithID(req.GetId())
	result, err := s.netboxClient.Dcim.DcimSitesRead(params, nil)
	if err != nil {
		return nil, err
	}
	site := result.Payload
	return &organization.Site{
		Id:          site.ID,
		Name:        *site.Name,
		Slug:        *site.Slug,
		Status:      swag.StringValue(site.Status.Value),
		Description: site.Description,
	}, nil
}

// TODO: Implementar CreateSite, UpdateSite, DeleteSite
