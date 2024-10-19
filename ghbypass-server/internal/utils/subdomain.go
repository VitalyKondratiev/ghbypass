package utils

import (
	"strings"
)

func GetSubdomain(host string, baseDomain string) string {
	host = strings.Split(host, ":")[0]

	if strings.HasSuffix(host, baseDomain) {
		subdomain := strings.TrimSuffix(host, "."+baseDomain)
		return subdomain
	}

	return ""
}
