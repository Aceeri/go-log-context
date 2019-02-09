package logContext

import (
	"net"
)

func LookupIp(dns string) string {
	addr, err := net.LookupIP(dns)
	if err == nil {
		for _, ip := range addr {
			if ip4 := ip.To4(); ip4 != nil {
				dns = ip4.String()
				break
			}
		}
	}

	return dns
}
