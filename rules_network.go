package validation

import (
	"net"
	"net/url"
	"regexp"
	"strings"
)

// Network validation rules

// UrlRule validates that a field is a valid URL
type UrlRule struct {
	Protocols []string
}

func (r *UrlRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	if str == "" {
		return false
	}
	
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	
	// Must have a scheme
	if u.Scheme == "" {
		return false
	}
	
	// If protocols are specified, check that the URL uses one of them
	if len(r.Protocols) > 0 {
		validProtocol := false
		for _, protocol := range r.Protocols {
			if u.Scheme == protocol {
				validProtocol = true
				break
			}
		}
		if !validProtocol {
			return false
		}
	}
	
	// Must have a host
	return u.Host != ""
}

func (r *UrlRule) Message() string {
	if len(r.Protocols) > 0 {
		return "The :attribute must be a valid URL with protocol: " + strings.Join(r.Protocols, ", ") + "."
	}
	return "The :attribute must be a valid URL."
}

// IpRule validates that a field is a valid IP address (IPv4 or IPv6)
type IpRule struct{}

func (r *IpRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	ip := net.ParseIP(str)
	return ip != nil
}

func (r *IpRule) Message() string {
	return "The :attribute must be a valid IP address."
}

// Ipv4Rule validates that a field is a valid IPv4 address
type Ipv4Rule struct{}

func (r *Ipv4Rule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	ip := net.ParseIP(str)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}

func (r *Ipv4Rule) Message() string {
	return "The :attribute must be a valid IPv4 address."
}

// Ipv6Rule validates that a field is a valid IPv6 address
type Ipv6Rule struct{}

func (r *Ipv6Rule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	ip := net.ParseIP(str)
	if ip == nil {
		return false
	}
	return ip.To4() == nil
}

func (r *Ipv6Rule) Message() string {
	return "The :attribute must be a valid IPv6 address."
}

// MacAddressRule validates that a field is a valid MAC address
type MacAddressRule struct{}

func (r *MacAddressRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	_, err := net.ParseMAC(str)
	return err == nil
}

func (r *MacAddressRule) Message() string {
	return "The :attribute must be a valid MAC address."
}

// HexColorRule validates that a field is a valid hexadecimal color
type HexColorRule struct{}

func (r *HexColorRule) Passes(attribute string, value interface{}) bool {
	str := ToString(value)
	// Regular expression for hex color (#RGB, #RRGGBB, #RRGGBBAA)
	hexColorRegex := regexp.MustCompile(`^#([A-Fa-f0-9]{3}|[A-Fa-f0-9]{6}|[A-Fa-f0-9]{8})$`)
	return hexColorRegex.MatchString(str)
}

func (r *HexColorRule) Message() string {
	return "The :attribute must be a valid hexadecimal color."
}