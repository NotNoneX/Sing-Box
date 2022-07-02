package option

import (
	"encoding/json"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

type RouteOptions struct {
	GeoIP *GeoIPOptions `json:"geoip,omitempty"`
	Rules []Rule        `json:"rules,omitempty"`
}

func (o RouteOptions) Equals(other RouteOptions) bool {
	return common.ComparablePtrEquals(o.GeoIP, other.GeoIP) &&
		common.SliceEquals(o.Rules, other.Rules)
}

type GeoIPOptions struct {
	Path           string `json:"path,omitempty"`
	DownloadURL    string `json:"download_url,omitempty"`
	DownloadDetour string `json:"download_detour,omitempty"`
}

type _Rule struct {
	Type           string       `json:"type,omitempty"`
	DefaultOptions *DefaultRule `json:"default_options,omitempty"`
	LogicalOptions *LogicalRule `json:"logical_options,omitempty"`
}

type Rule _Rule

func (r Rule) Equals(other Rule) bool {
	return r.Type == other.Type &&
		common.PtrEquals(r.DefaultOptions, other.DefaultOptions) &&
		common.PtrEquals(r.LogicalOptions, other.LogicalOptions)
}

func (r *Rule) MarshalJSON() ([]byte, error) {
	var v any
	switch r.Type {
	case C.RuleTypeDefault:
		v = r.DefaultOptions
	case C.RuleTypeLogical:
		v = r.LogicalOptions
	default:
		return nil, E.New("unknown rule type: " + r.Type)
	}
	return MarshallObjects(r, v)
}

func (r *Rule) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_Rule)(r))
	if err != nil {
		return err
	}
	if r.Type == "" {
		r.Type = C.RuleTypeDefault
	}
	var v any
	switch r.Type {
	case C.RuleTypeDefault:
		v = &r.DefaultOptions
	case C.RuleTypeLogical:
		v = &r.LogicalOptions
	default:
		return E.New("unknown rule type: " + r.Type)
	}
	return json.Unmarshal(bytes, v)
}

type DefaultRule struct {
	Inbound       Listable[string] `json:"inbound,omitempty"`
	IPVersion     int              `json:"ip_version,omitempty"`
	Network       string           `json:"network,omitempty"`
	Protocol      Listable[string] `json:"protocol,omitempty"`
	Domain        Listable[string] `json:"domain,omitempty"`
	DomainSuffix  Listable[string] `json:"domain_suffix,omitempty"`
	DomainKeyword Listable[string] `json:"domain_keyword,omitempty"`
	SourceGeoIP   Listable[string] `json:"source_geoip,omitempty"`
	GeoIP         Listable[string] `json:"geoip,omitempty"`
	SourceIPCIDR  Listable[string] `json:"source_ip_cidr,omitempty"`
	IPCIDR        Listable[string] `json:"ip_cidr,omitempty"`
	SourcePort    Listable[uint16] `json:"source_port,omitempty"`
	Port          Listable[uint16] `json:"port,omitempty"`
	// ProcessName   Listable[string] `json:"process_name,omitempty"`
	// ProcessPath   Listable[string] `json:"process_path,omitempty"`
	Outbound string `json:"outbound,omitempty"`
}

func (r DefaultRule) IsValid() bool {
	var defaultValue DefaultRule
	defaultValue.Outbound = r.Outbound
	return !r.Equals(defaultValue)
}

func (r DefaultRule) Equals(other DefaultRule) bool {
	return common.ComparableSliceEquals(r.Inbound, other.Inbound) &&
		r.IPVersion == other.IPVersion &&
		r.Network == other.Network &&
		common.ComparableSliceEquals(r.Protocol, other.Protocol) &&
		common.ComparableSliceEquals(r.Domain, other.Domain) &&
		common.ComparableSliceEquals(r.DomainSuffix, other.DomainSuffix) &&
		common.ComparableSliceEquals(r.DomainKeyword, other.DomainKeyword) &&
		common.ComparableSliceEquals(r.SourceGeoIP, other.SourceGeoIP) &&
		common.ComparableSliceEquals(r.GeoIP, other.GeoIP) &&
		common.ComparableSliceEquals(r.SourceIPCIDR, other.SourceIPCIDR) &&
		common.ComparableSliceEquals(r.IPCIDR, other.IPCIDR) &&
		common.ComparableSliceEquals(r.SourcePort, other.SourcePort) &&
		common.ComparableSliceEquals(r.Port, other.Port) &&
		r.Outbound == other.Outbound
}

type LogicalRule struct {
	Mode     string        `json:"mode"`
	Rules    []DefaultRule `json:"rules,omitempty"`
	Outbound string        `json:"outbound,omitempty"`
}

func (r LogicalRule) IsValid() bool {
	return len(r.Rules) > 0 && common.All(r.Rules, DefaultRule.IsValid)
}

func (r LogicalRule) Equals(other LogicalRule) bool {
	return r.Mode == other.Mode &&
		common.SliceEquals(r.Rules, other.Rules) &&
		r.Outbound == other.Outbound
}