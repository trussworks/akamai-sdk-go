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

// ZoneDeleteRequest is a slice of zones to delete when making a request to the Akamai API.
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

// RecordSet is set of DNS records belonging to a particular DNS name
type RecordSet struct {
	Name  *string   `json:"name,omitempty"`
	Rdata []*string `json:"rdata,omitempty"`
	TTL   *int      `json:"ttl,omitempty"`
	Type  *string   `json:"type,omitempty"`
}

// RecordSetOptions specifies optional parameters to the FastDNSv2ServiceListRecordSet method.
type RecordSetOptions struct {
	Zone string `json:"zone,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// ListRecordSet retrieves a single record set for the zone, record name, and record type specified in the URL.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getzonerecordset
func (s *FastDNSv2Service) ListRecordSet(ctx context.Context, opt *RecordSetOptions) (*RecordSet, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/zones/%v/names/%v/types/%v", opt.Zone, opt.Name, opt.Type)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var rs *RecordSet
	resp, err := s.client.Do(ctx, req, &rs)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, nil
}

// CreateRecordSet creates a new Record Set with the specified name and type.
//
// https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#postzonerecordset
func (s *FastDNSv2Service) CreateRecordSet(ctx context.Context, opt *RecordSetOptions, rs *RecordSet) (*RecordSet, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/zones/%v/names/%v/types/%v", opt.Zone, opt.Name, opt.Type)

	req, err := s.client.NewRequest("POST", u, rs)
	if err != nil {
		return nil, nil, err
	}

	var r *RecordSet
	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// UpdateRecordSet replaces an existing Record Set with the request body.
// The name and type must match the existing record.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#putzonerecordset
func (s *FastDNSv2Service) UpdateRecordSet(ctx context.Context, opt *RecordSetOptions, rs *RecordSet) (*RecordSet, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/zones/%v/names/%v/types/%v", opt.Zone, opt.Name, opt.Type)

	req, err := s.client.NewRequest("PUT", u, rs)
	if err != nil {
		return nil, nil, err
	}

	var r *RecordSet
	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// DeleteRecordSet removes an existing record set.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#deletezonerecordset
func (s *FastDNSv2Service) DeleteRecordSet(ctx context.Context, opt *RecordSetOptions) (*Response, error) {
	u := fmt.Sprintf("/config-dns/v2/zones/%v/names/%v/types/%v", opt.Zone, opt.Name, opt.Type)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
