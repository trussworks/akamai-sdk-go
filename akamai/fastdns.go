package akamai

import (
	"context"
	"fmt"
)

// FastDNSv2Service handles communication with the v2 FastDNS (beta) related endpoints
// of the Akamai API
type FastDNSv2Service service

// Zone represents an Akamai zone from the v2 FastDNS API.
type Zone struct {
	ContractID         *string   `json:"contractId,omitempty"`
	Zone               *string   `json:"zone,omitempty"`
	Type               *string   `json:"type,omitempty"`
	Comment            *string   `json:"comment,omitempty"`
	EndCustomerID      *string   `json:"endCustomerId,omitempty"`
	Target             *string   `json:"target,omitempty"`
	TSIGKey            *tsigKey  `json:"tsigKey,omitempty"`
	Masters            []*string `json:"masters,omitempty"`
	VersionID          *string   `json:"versionId,omitempty"`
	LastModifiedDate   *string   `json:"lastModifiedDate,omitempty"`
	LastModifiedBy     *string   `json:"lastModifiedBy,omitempty"`
	LastActivationDate *string   `json:"lastActivationDate,omitempty"`
	ActivationState    *string   `json:"activationState,omitempty"`
}

type tsigKey struct {
	Name      *string `json:"name,omitempty"`
	Algorithm *string `json:"algorithm,omitempty"`
	Secret    *string `json:"secret,omitempty"`
}

// ZoneListOptions specifies optional parameters to the FastDNSv2Service.ListZones method.
type ZoneListOptions struct {
	ContractIDs string `url:"contractIds,omitempty"`
	Page        int    `url:"page,omitempty"`
	PageSize    int    `url:"pageSize,omitempty"`
	Search      string `url:"search,omitempty"`
	ShowAll     bool   `url:"showAll,omitempty"`
	SortBy      string `url:"sortBy,omitempty"`
	Types       string `url:"types,omitempty"`
	GroupID     int    `url:"gid,omitempty"`
}

// ListZones retreives the zones for the authenticated user.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getzones
func (s *FastDNSv2Service) ListZones(ctx context.Context, opt *ZoneListOptions) ([]*Zone, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones")
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var zones []*Zone
	resp, err := s.client.Do(ctx, req, &zones)
	if err != nil {
		return nil, resp, err
	}

	return zones, resp, nil
}

// CreateZone creates a new Zone
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#postzones
func (s *FastDNSv2Service) CreateZone(ctx context.Context, opt *ZoneListOptions, zone *Zone) (*Zone, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones")
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, zone)
	if err != nil {
		return nil, nil, err
	}

	z := new(Zone)
	resp, err := s.client.Do(ctx, req, z)
	if err != nil {
		return nil, resp, err
	}

	return z, resp, nil
}

// UpdateZone modifies an Akamai zone.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#putzone
func (s *FastDNSv2Service) UpdateZone(ctx context.Context, zone *Zone) (*Zone, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones/%v", zone.Zone)
	req, err := s.client.NewRequest("PUT", u, zone)

	if err != nil {
		return nil, nil, err
	}

	z := new(Zone)
	resp, err := s.client.Do(ctx, req, z)
	if err != nil {
		return nil, resp, err
	}

	return z, resp, nil
}

type ZoneDeleteRequest struct {
	Zones []*string `json:"zones,omitempty"`
}

// DeleteZone deletes one or more Akamai zones.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#postbulkzonedelete
func (s *FastDNSv2Service) DeleteZone(ctx context.Context, zd *ZoneDeleteRequest) (*Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones/delete-requests")
	req, err := s.client.NewRequest("POST", u, zd)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
