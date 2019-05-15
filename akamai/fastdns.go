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

// ZoneCreateOptions specifies the optional parameters to the FastDNSV2Service.CreateZone method.
type ZoneCreateOptions struct {
	ContractID string `url:"contractId,omitempty"`
}

//ZoneCreateRequest specifies the parameters for the CreateZone method.
type ZoneCreateRequest struct {
	Zone             string   `json:"zone,omitempty"`
	Type             string   `json:"type,omitempty"`
	Comment          string   `json:"comment,omitempty"`
	EndCustomerID    string   `json:"endCustomerId,omitempty"`
	Target           string   `json:"target,omitempty"`
	TSIGKey          string   `json:"tsigKey,omitempty"`
	Masters          []string `json:"masters,omitempty"`
	SignAndServe     bool     `json:"signAndServe"`
	SignAndServeAlgo string   `json:"signAndServeAlgorithm,omitempty"`
}

// ZoneList holds a response from ListZones
type ZoneList struct {
	Metadata *ZoneListMetadata `json:"metadata,omitempty"`
	Zones    []*Zone           `json:"zones,omitempty"`
}

// ZoneListMetadata holds metadata from the ZoneList response
type ZoneListMetadata struct {
	ContractIDs   []*string `json:"contractId,omitempty"`
	Page          *int      `json:"page,omitempty"`
	PageSize      *int      `json:"pageSize,omitempty"`
	ShowAll       *bool     `json:"showAll,omitempty"`
	TotalElements *int      `json:"totalElements,omitempty"`
}

// ZoneMetadata holds the response from GetZone
type ZoneMetadata struct {
	ContractID            *string `json:"contractId,omitempty"`
	Zone                  *string `json:"zone,omitempty"`
	Type                  *string `json:"type,omitempty"`
	AliasCount            *int    `json:"aliasCount,omitempty"`
	SignAndServe          *bool   `json:"signAndServe,omitempty"`
	SignAndServeAlgorithm *string `json:"signAndServeAlgorithm,omitempty"`
	VersionId             *string `json:"versionId,omitempty"`
	LastModifiedDate      *string `json:"lastModifiedDate,omitempty"`
	LastModifiedBy        *string `json:"lastModifiedBy,omitempty"`
	LastActivationDate    *string `json:"lastActivationDate,omitempty"`
	ActivationState       *string `json:"activationState,omitempty"`
	Comment               *string `json:"comment,omitempty"`
}

// ListZones retreives the zones for the authenticated user.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getzones
func (s *FastDNSv2Service) ListZones(ctx context.Context, opt *ZoneListOptions) (*ZoneList, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones")
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var zones *ZoneList
	resp, err := s.client.Do(ctx, req, &zones)
	if err != nil {
		return nil, resp, err
	}

	return zones, resp, nil
}

// GetZone retrieves the metadata of a single zone. Does not include record sets.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getzone
func (s *FastDNSv2Service) GetZone(ctx context.Context, zone string) (*ZoneMetadata, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones/%v", zone)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var zmeta *ZoneMetadata
	resp, err := s.client.Do(ctx, req, &zmeta)
	if err != nil {
		return nil, resp, err
	}

	return zmeta, resp, nil
}

// CreateZone creates a new Zone
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#postzones
func (s *FastDNSv2Service) CreateZone(ctx context.Context, cid string, zone *ZoneCreateRequest) (*Zone, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones")
	lo := ZoneCreateOptions{
		ContractID: cid,
	}

	u, err := addOptions(u, lo)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, zone)
	if err != nil {
		return nil, nil, err
	}

	z := new(Zone)
	resp, err := s.client.Do(ctx, req, &z)
	if err != nil {
		return nil, resp, err
	}

	return z, resp, nil
}

// UpdateZone modifies an Akamai zone.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#putzone
func (s *FastDNSv2Service) UpdateZone(ctx context.Context, zone *ZoneCreateRequest) (*Zone, *Response, error) {
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
	Zones []string `json:"zones,omitempty"`
}

// ZoneDeleteResponse holds the response we must poll to confirm request success/failure.
type ZoneDeleteResponse struct {
	RequestID      *string `json:"requestId,omitempty"`
	ExpirationDate *string `json:"expirationDate,omitempty"`
	ZonesSubmitted *int    `json:"zonesSubmitted,omitempty"`
	SuccessCount   *int    `json:"successCount,omitempty"`
	FailureCount   *int    `json:"failureCount,omitempty"`
	IsComplete     *bool   `json:"isComplete"`
}

// ZoneDeleteResult holds the result of  the ZoneDelete request
type ZoneDeleteResult struct {
	RequestID    *string   `json:"requestId,omitempty"`
	DeletedZones []*string `json:"successfullyDeletedZones,omitempty"`
	FailedZones  []*struct {
		Zone          *string `json:"zone,omitempty"`
		FailureReason *string `json:"failiureReason,omitempty"`
	} `json:"failedZones,omitempty"`
}

// DeleteZone deletes one or more Akamai zones.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#postbulkzonedelete
func (s *FastDNSv2Service) DeleteZone(ctx context.Context, zd *ZoneDeleteRequest) (*ZoneDeleteResponse, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones/delete-requests")
	req, err := s.client.NewRequest("POST", u, zd)
	if err != nil {
		return nil, nil, err
	}

	z := new(ZoneDeleteResponse)
	resp, err := s.client.Do(ctx, req, &z)
	if err != nil {
		return nil, resp, err
	}

	return z, resp, nil
}

// DeleteZoneStatus checks the status of a DeleteZone request. Use the request ID that was given
// from teh output of the DeleteZone request.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getbulkzonedeletestatus
func (s *FastDNSv2Service) DeleteZoneStatus(ctx context.Context, rid string) (*ZoneDeleteResponse, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones/delete-requests/%v", rid)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	z := new(ZoneDeleteResponse)
	resp, err := s.client.Do(ctx, req, &z)
	if err != nil {
		return nil, resp, err
	}

	return z, resp, nil
}

// DeleteZoneResult retrieves the results from a completed DeleteZone request.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getbulkzonedeleteresult
func (s *FastDNSv2Service) DeleteZoneResult(ctx context.Context, rid string) (*ZoneDeleteResult, *Response, error) {
	u := fmt.Sprintf("config-dns/v2/zones/delete-requests/%v/result", rid)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	z := new(ZoneDeleteResult)
	resp, err := s.client.Do(ctx, req, &z)
	if err != nil {
		return nil, resp, err

	}

	return z, resp, nil
}

// RecordSet is set of DNS records belonging to a particular DNS name
type RecordSet struct {
	Name  *string   `json:"name,omitempty"`
	Rdata []*string `json:"rdata,omitempty"`
	TTL   *int      `json:"ttl,omitempty"`
	Type  *string   `json:"type,omitempty"`
	State *string   `json:"state,omitempty"`
}

// RecordSetOptions specifies optional parameters to some record set methods.
type RecordSetOptions struct {
	Zone string `json:"zone,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// RecordSetCreatRequest is a set of DNS records belonging to a particular DNS name to be created.
type RecordSetCreateRequest struct {
	Zone  string   `json:"zone,omitempty"`
	Name  string   `json:"name,omitempty"`
	Rdata []string `json:"rdata,omitempty"`
	TTL   int      `json:"ttl,omitempty"`
	Type  string   `json:"type,omitempty"`
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
func (s *FastDNSv2Service) CreateRecordSet(ctx context.Context, rs *RecordSetCreateRequest) (*RecordSet, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/zones/%v/names/%v/types/%v", rs.Zone, rs.Name, rs.Type)

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
func (s *FastDNSv2Service) UpdateRecordSet(ctx context.Context, rs *RecordSetCreateRequest) (*RecordSet, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/zones/%v/names/%v/types/%v", rs.Zone, rs.Name, rs.Type)

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

// ListZoneRecordSets holds the response from GetZoneRecordSets.
type ListZoneRecordSets struct {
	Metadata   *ListZoneRecordMetadata `json:"metadata,omitempty"`
	RecordSets []*RecordSet            `json:"recordsets,omitempty"`
}

// ListZoneRecordMetadata holds the metadata response from GetZoneRecordSets.
type ListZoneRecordMetadata struct {
	Zone          *string   `json:"zone,omitempty"`
	Types         []*string `json:"types,omitempty"`
	Page          *int      `json:"page,omitempty"`
	PageSize      *int      `json:"pageSize,omitempty"`
	TotalElements *int      `json:"totalElements,omitempty"`
}

// ListZoneRecordSetOptions are optional query parameters.
type ListZoneRecordSetOptions struct {
	Page     int    `url:"page,omitempty"`
	PageSize int    `url:"pageSize,omitempty"`
	Search   string `url:"search,omitempty"`
	ShowAll  bool   `url:"showAll,omitempty"`
	SortBy   string `url:"sortBy,omitempty"`
	Types    string `url:"types,omitempty"`
}

// GetZoneRecordSets lists all record sets for this zone. Can only be used on PRIMARY
// and SECONDARY zones. This operation is paginated.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getzonerecordsets
func (s *FastDNSv2Service) GetZoneRecordSets(ctx context.Context, zone string, opt *ListZoneRecordSetOptions) (*ListZoneRecordSets, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/zones/%v/recordsets", zone)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var z *ListZoneRecordSets
	resp, err := s.client.Do(ctx, req, &z)
	if err != nil {
		return nil, resp, err
	}

	return z, resp, nil
}

// Contract holds Akamai's Contract object type. It provides metadata about
// a customer's Akamai FastDNS account.
type Contract struct {
	ContractID       *string   `json:"contractId,omitempty"`
	ContractName     *string   `json:"contractName,omitempty"`
	ContractTypeName *string   `json:"contractTypeName,omitempty"`
	Features         []*string `json:"features,omitempty"`
	Permissions      []*string `json:"permissions,omitempty"`
	ZoneCount        int       `json:"zoneCount,omitempty"`
	MaximumZones     int       `json:"maximumZones,omitempty"`
}

// GetZoneContract returns data about the Contract to which the Zone belongs.
//
// Akamai API docs: https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getzonecontract
func (s *FastDNSv2Service) GetZoneContract(ctx context.Context, zone string) (*Contract, *Response, error) {

	u := fmt.Sprintf("/config-dns/v2/zones/%v/contract", zone)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var c *Contract
	resp, err := s.client.Do(ctx, req, &s)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// ChangeListOptions holds options to pass when creating change lists.
type ChangeListOptions struct {
	Zone      string `url:"zone,omitempty"`
	Overwrite string `url:"overwrite,omitempty"`
	Page      int    `url:"page,omitempty"`
	PageSize  int    `url:"pageSize,omitempty"`
	Search    string `url:"search,omitempty"`
	ShowAll   bool   `url:"showAll,omitempty"`
	SortBy    string `url:"sortBy,omitempty"`
	Types     string `url:"types,omitempty"`
}

// ChangeList holds metadata about a change list, including the particular version of a zone that
// the change list was based off when it was created.
type ChangeList struct {
	ChangeTag        string `json:"changeTag,omitempty"`
	LastModifiedDate string `json:"lastModifiedDate,omitempty"`
	Stale            string `json:"stale,omitempty"`
	Zone             string `json:"zone,omitempty"`
	ZoneVersionId    string `json:"zoneVersionId,omitempty"`
}

// CreateChangeList creates a new Change List based on the most recent version of a zone.
// No POST body is required, since the object is read-only.
//
// Akamai API docs:
// https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#postchangelists
func (s *FastDNSv2Service) CreateChangeList(ctx context.Context, cl *ChangeListOptions) (*ChangeList, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/changelists")
	u, err := addOptions(u, cl)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, nil)

	c := new(ChangeList)
	resp, err := s.client.Do(ctx, req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// GetChangeList describes a Change List, showing its base zone version,
// last modified time, and current change tag.
//
// Akamai API docs:
// https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getchangelist
func (s *FastDNSv2Service) GetChangeList(ctx context.Context, zone string) (*ChangeList, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/changelists/%v", zone)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	c := new(ChangeList)
	resp, err := s.client.Do(ctx, req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// ChangeListRecords holds the current list of record sets from the perspective of a change list
type ChangeListRecords struct {
	Metadata   *ChangeListMetadata `json:"metadata,omitempty"`
	Recordsets []*RecordSet        `json:"recordsets,omitempty"`
}

// ChangeListMetadata holds metadata for a change list of record sets
type ChangeListMetadata struct {
	Zone          *string   `json:"zone,omitempty"`
	Types         []*string `json:"types,omitempty"`
	Page          *int      `json:"page,omitempty"`
	PageSize      *int      `json:"pageSize,omitempty"`
	TotalElements *int      `json:"totalElements,omitempty"`
}

// GetChangeListRecordSets retrieves the current list of record sets from the perspective of this change list.
// Any changes that have been added to this change list will be reflected in the list of record sets returned.
//
// Akamai API docs:
// https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#getchangelistrecordsets
func (s *FastDNSv2Service) GetChangeListRecordSets(ctx context.Context, zone string, opt *ChangeListOptions) (*ChangeListRecords, *Response, error) {
	u := fmt.Sprintf("/config-dns/v2/changelists/%v/recordsets", zone)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	c := new(ChangeListRecords)
	resp, err := s.client.Do(ctx, req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil

}

// DeleteChangeList removes an unneeded Change List
//
// Akamai API docs:
// https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#deletechangelist
func (s *FastDNSv2Service) DeleteChangeList(ctx context.Context, zone string) (*Response, error) {
	u := fmt.Sprintf("/config-dns/v2/changelists/%v", zone)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// SubmitChangeList applies all of the changes in this change list to the current zone. This
// operation fails if the change list has become stale.
//
// Akamai API docs:
// https://developer.akamai.com/api/web_performance/fast_dns_zone_management/v2.html#postchangelistsubmit
func (s *FastDNSv2Service) SubmitChangeList(ctx context.Context, zone string) (*Response, error) {
	u := fmt.Sprintf("/config-dns/v2/changelists/%v/submit", zone)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
