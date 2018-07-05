package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"xframe/server/websocket"
)

//ratelimit
var (
	defaultStrategy string  = "token"
	defaultRate     float64 = 0.0
	defaultCapacity int64   = 0
)

var (
	OnDataIn     = func(conn *TcpConnection, msg []byte) {}
	OnDataOut    = func(conn net.Conn, msg []byte) {}
	OnConnect    = func(conn *TcpConnection) {}
	OnDisconnect = func(conn *TcpConnection) {}
	RouteHTTP    = func(w http.ResponseWriter, r *http.Request) {}
	RouteWs      = func(ws *websocket.Conn) {}
)

func isIPv4(ip string) bool {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return false
	}
	return true
}

func isIPv6(ip string) bool {
	return strings.Contains(ip, ":")
}

func isIP(ip string) bool {
	addr := net.ParseIP(ip)
	if addr == nil {
		return false
	}
	return true
}

func parseListenAddr(addr string) (ip string, err error) {
	if len(addr) <= 0 {
		err = errors.New("addr is empty")
		return
	}
	if addr[0] == '@' {
		if ip, err = getAddrByNetIf(string([]byte(addr)[1:])); err != nil {
			return
		}
	} else if isIP(addr) {
		ip = addr
	} else {
		err = errors.New(fmt.Sprintf("parse [\"%s\"] is not a ip", addr))
	}
	return

}

func getAddrByNetIf(netIf string) (addr string, err error) {
	net_fields := strings.Split(netIf, ":")
	nic := net_fields[0]
	net_flag := "ipv4"
	if len(net_fields) > 1 {
		net_flag = net_fields[1]
	}
	ifi, err := net.InterfaceByName(nic)
	if err != nil {
		return
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		return
	}
	for _, a := range addrs {
		ip, _, e := net.ParseCIDR(a.String())
		if e != nil {
			err = e
			return
		}
		if net_flag == "ipv4" {
			if ipv4 := ip.To4(); isIPv4(ipv4.String()) {
				addr = ipv4.String()
				return
			}
		} else if net_flag == "ipv6" {
			if ipv6 := ip.To16(); isIPv6(ipv6.String()) {
				addr = ipv6.String()
				return
			}
		}
	}
	err = errors.New(fmt.Sprintf("can't get interface[\" %s \"]'s ip", netIf))
	return
}
