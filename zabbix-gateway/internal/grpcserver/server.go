package grpcserver

import (
	"context"
	"zabbix-gateway/internal/zabbix_client"
	"zabbix-gateway/proto/zabbix"
)

type Server struct {
	monitoring.UnimplementedMonitoringServiceServer
	zabbixClient *zabbix_client.Client
}

func NewServer(zabbixClient *zabbix_client.Client) *Server {
	return &Server{zabbixClient: zabbixClient}
}

func (s *Server) ListHostGroups(ctx context.Context, req *monitoring.ListHostGroupsRequest) (*monitoring.ListHostGroupsResponse, error) {
	groups, err := s.zabbixClient.ListHostGroups(ctx)
	if err != nil {
		return nil, err
	}
	protoGroups := make([]*monitoring.HostGroup, len(groups))
	for i, g := range groups {
		protoGroups[i] = &monitoring.HostGroup{Groupid: g.ID, Name: g.Name}
	}
	return &monitoring.ListHostGroupsResponse{Groups: protoGroups}, nil
}

func (s *Server) ListHosts(ctx context.Context, req *monitoring.ListHostsRequest) (*monitoring.ListHostsResponse, error) {
	hosts, err := s.zabbixClient.ListHostsByGroupID(ctx, req.GetGroupids())
	if err != nil {
		return nil, err
	}
	protoHosts := make([]*monitoring.Host, len(hosts))
	for i, h := range hosts {
		protoHosts[i] = &monitoring.Host{Hostid: h.ID, Host: h.Host, Name: h.Name}
	}
	return &monitoring.ListHostsResponse{Hosts: protoHosts}, nil
}
func (s *Server) ListItems(ctx context.Context, req *monitoring.ListItemsRequest) (*monitoring.ListItemsResponse, error) {
	items, err := s.zabbixClient.ListItemsByHostID(ctx, req.GetHostids())
	if err != nil {
		return nil, err
	}
	protoItems := make([]*monitoring.Item, len(items))
	for i, item := range items {
		protoItems[i] = &monitoring.Item{
			Itemid:    item.ID,
			Name:      item.Name,
			Key_:      item.Key_,
			Lastvalue: item.LastValue,
			Lastclock: item.LastClock,
		}
	}
	return &monitoring.ListItemsResponse{Items: protoItems}, nil
}
func (s *Server) ListAlerts(ctx context.Context, req *monitoring.ListAlertsRequest) (*monitoring.ListAlertsResponse, error) {
	alerts, err := s.zabbixClient.ListRecentAlertsByHostID(ctx, req.GetHostids())
	if err != nil {
		return nil, err
	}
	protoAlerts := make([]*monitoring.Alert, len(alerts))
	for i, alert := range alerts {
		protoAlerts[i] = &monitoring.Alert{
			Triggerid:   alert.TriggerID,
			Description: alert.Description,
			Priority:    alert.Priority,
			Lastchange:  alert.LastChange,
			Value:       alert.Value,
		}
	}
	return &monitoring.ListAlertsResponse{Alerts: protoAlerts}, nil
}
