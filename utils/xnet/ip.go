package xnet

import (
	"net"
	"os"
	"strings"
)

// IsPublicIP 判断ip是否是公网IP
func IsPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

// IsIPv4 判断是否是IPv4
func IsIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

// GetLanIPs 获取局域网本地IP地址列表
func GetLanIPs() (ips []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			!ipnet.IP.IsLinkLocalMulticast() &&
			!ipnet.IP.IsLinkLocalUnicast() &&
			IsIPv4(ipnet.IP) &&
			!IsPublicIP(ipnet.IP) {

			ips = append(ips, ipnet.IP.String())
		}
	}

	return
}

func HostName() string {
	host, _ := os.Hostname()
	return host
}

func ServerIPString() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	ips := make([]string, 0)
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}

		}
	}

	return strings.Join(ips, ",")
}

func ServerIP() (string, string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", ""
	}
	var internalIP, externalIP string
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipaddr := ipnet.IP.String()
				if strings.HasPrefix(ipaddr, "10.") || strings.HasPrefix(ipaddr, "192.") || strings.HasPrefix(ipaddr, "127.") {
					internalIP = ipaddr
				} else {
					externalIP = ipaddr
				}
			}

		}
	}

	return externalIP, internalIP
}

func SerIP() string {
	eIP, iIP := ServerIP()
	if eIP != "" {
		return eIP
	} else {
		return iIP
	}
}
