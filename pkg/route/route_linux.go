//go:build linux
// +build linux

package route

import (
	"fmt"
	"net"
	"strings"
	"syscall"
	"unsafe"
)

func Add(route Route) error {
	return nil
}

func Delete(route Route) error {
	return nil
}

type routeInfoInMemory struct {
	Family   byte
	DstLen   byte
	SrcLen   byte
	TOS      byte
	Table    byte
	Protocol byte
	Scope    byte
	Type     byte
	Flags    uint32
}

type RouteRule struct {
	Src         *net.IPNet
	Dst         *net.IPNet
	Gateway     net.IP
	PrefSrc     net.IP
	InputIface  uint32
	OutputIface uint32
	Priority    uint32
}

type routeManager struct {
	ifaces map[int]net.Interface
	addrs  map[int]ipAddrs
	v4     []Route
}

type ipAddrs struct {
	v4 net.IP
	v6 net.IP
}

func (manager *routeManager) newRouteRule(rt *routeInfoInMemory, attrs []syscall.NetlinkRouteAttr) *Route {
	ret := &Route{}
	for _, attr := range attrs {
		switch attr.Attr.Type {
		case syscall.RTA_DST:
			ret.Destination = net.IP(attr.Value).String()
			ret.Mask = net.IP(net.CIDRMask(int(rt.DstLen), len(attr.Value)*8)).String()
		case syscall.RTA_SRC:
			// ret.Src = &net.IPNet{
			// 	IP:   net.IP(attr.Value),
			// 	Mask: net.CIDRMask(int(rt.SrcLen), len(attr.Value)*8),
			// }
		case syscall.RTA_GATEWAY:
			ret.Gateway = net.IP(attr.Value).String()
		case syscall.RTA_PREFSRC:
			// ret.PrefSrc = net.IP(attr.Value)
		case syscall.RTA_IIF:
			// ret.InputIface = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
		case syscall.RTA_OIF:
			inetIndex := *(*uint32)(unsafe.Pointer(&attr.Value[0]))
			if iface, ok := manager.ifaces[int(inetIndex)]; ok {
				ret.InterfaceName = iface.Name
			}

			if addr, ok := manager.addrs[int(inetIndex)]; ok {
				ret.InterfaceIP = addr.v4.String()
			}
			// ret.OutputIface = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
		case syscall.RTA_PRIORITY:
			// ret.Priority = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
		default:
			// log.Printf("not match %d", attr.Attr.Type)
		}
	}
	return ret
}

func (rtr *routeManager) initInterfaces() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		rtr.ifaces[iface.Index] = iface
		var addrs ipAddrs
		ifaceAddrs, err := iface.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range ifaceAddrs {
			if inet, ok := addr.(*net.IPNet); ok {
				if v4 := inet.IP.To4(); v4 != nil {
					if addrs.v4 == nil {
						addrs.v4 = v4
					}
				} else if addrs.v6 == nil {
					addrs.v6 = inet.IP
				}
			}
		}
		rtr.addrs[iface.Index] = addrs
	}
	return nil
}

func List() ([]Route, error) {
	rtr := &routeManager{}
	rtr.ifaces = make(map[int]net.Interface)
	rtr.addrs = make(map[int]ipAddrs)
	if err := rtr.initInterfaces(); err != nil {
		return nil, err
	}
	tab, err := syscall.NetlinkRIB(syscall.RTM_GETROUTE, syscall.AF_UNSPEC)
	if err != nil {
		return nil, err
	}
	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, err
	}
loop:
	for _, m := range msgs {
		switch m.Header.Type {
		case syscall.NLMSG_DONE:
			break loop
		case syscall.RTM_NEWROUTE:
			rt := (*routeInfoInMemory)(unsafe.Pointer(&m.Data[0]))
			attrs, err := syscall.ParseNetlinkRouteAttr(&m)
			if err != nil {
				return nil, err
			}
			routeInfo := rtr.newRouteRule(rt, attrs)
			if rt.Family != syscall.AF_INET && rt.Family != syscall.AF_INET6 {
				continue
			}

			if routeInfo.Destination == "" && routeInfo.Gateway == "" {
				continue
			}
			switch rt.Family {
			case syscall.AF_INET:
				rtr.v4 = append(rtr.v4, *routeInfo)
			default:
				continue
			}
		}
	}

	return rtr.v4, nil
}

func (r *routeManager) String() string {
	strs := []string{}
	for _, route := range r.v4 {
		strs = append(strs, fmt.Sprintf("dst %s mask %s interfaceName %s interfaceIP %s",
			route.Destination, route.Mask, route.InterfaceName, route.InterfaceIP))
	}
	return strings.Join(strs, "\n")
}
