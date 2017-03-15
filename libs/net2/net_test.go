package net2

import (
	"fmt"
	"net"
	"testing"
)

func TestParseHosts(t *testing.T) {
	hosts := "127.0.0.1:2379;192.168.174.114:2379"
	arr := ParseHosts(hosts)
	t.Log(arr)
}

func TestCheckIP(t *testing.T) {
	ipStr := "192.168.174.114:2379"
	ip := CheckIp(ipStr)
	t.Log(ip)
}

func TestGetIp(t *testing.T) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(err)
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(err)
			continue
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					fmt.Printf("%v : %s (%s)\n", i.Name, ipnet.IP.String(), ipnet.IP.DefaultMask())
				}
			}
			// //fmt.Println(a)
			// switch v := a.(type) {
			// case *net.IPNet:
			// 	fmt.Printf("%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())
			// }

		}
	}
}
