package main

import (
	"net"
	"net/http"
	"strings"
)

// list of private subnets
var privateMasks, _ = toMasks([]string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"fc00::/7",
})

// converts a list of subnets' string to a list of net.IPNet.
func toMasks(ips []string) (masks []net.IPNet, err error) {
	for _, cidr := range ips {
		var network *net.IPNet
		_, network, err = net.ParseCIDR(cidr)
		if err != nil {
			return
		}
		masks = append(masks, *network)
	}
	return
}

// checks if a net.IP is in a list of net.IPNet
func ipInMasks(ip net.IP, masks []net.IPNet) bool {
	for _, mask := range masks {
		if mask.Contains(ip) {
			return true
		}
	}
	return false
}

// IsPublicIP returns true if the given IP can be routed on the Internet.
func IsPublicIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() {
		return false
	}
	return !ipInMasks(ip, privateMasks)
}

// Parse parses the value of the X-Forwarded-For Header and returns the IP address.
func Parse(xffString string, numProxies int) string {
	ipList := strings.Split(xffString, ",")
	if numProxies <= 0 {
		numProxies = len(ipList)
	}
	for idx := len(ipList) - 1; idx >= 0; idx-- {
		ipStr := strings.TrimSpace(ipList[idx])
		if ip := net.ParseIP(ipStr); ip != nil && IsPublicIP(ip) {
			if numProxies <= 0 || idx == 0 {
				return ipStr
			}
			numProxies--
		}
	}
	return ""
}

// GetRemoteAddr parses the given request, resolves the X-Forwarded-For header
// and returns the resolved remote address.
func GetRemoteAddr(r *http.Request, numProxies int) string {
	if xffh := r.Header.Get("X-Forwarded-For"); xffh != "" {
		if sip, sport, err := net.SplitHostPort(r.RemoteAddr); err == nil && sip != "" {
			if xip := Parse(xffh, numProxies); xip != "" {
				return net.JoinHostPort(xip, sport)
			}
		}
	}
	return r.RemoteAddr
}
