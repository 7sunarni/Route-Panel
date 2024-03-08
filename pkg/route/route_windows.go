package route

import (
	"errors"
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

const routeMetric = 93

type routeManager struct {
	iphlpapi             *syscall.LazyDLL
	getIpForwardTable    *syscall.LazyProc
	createIpForwardEntry *syscall.LazyProc
	deleteIpForwardEntry *syscall.LazyProc
	netInterfaces        map[int]netInterface
}

type routeRule struct {
	Dest      [4]byte //目标网络
	Mask      [4]byte //掩码
	Policy    uint32  //ForwardPolicy:0x0
	NextHop   [4]byte //网关
	IfIndex   uint32  // 网卡索引 id
	Type      uint32  //3 本地接口  4 远端接口
	Proto     uint32  //3静态路由 2本地接口 5EGP网关
	Age       uint32  //存在时间 秒
	NextHopAS uint32  //下一跳自治域号码 0
	Metric1   uint32  //度量衡(跃点数)，根据 ForwardProto 不同意义不同。
	Metric2   uint32
	Metric3   uint32
	Metric4   uint32
	Metric5   uint32
}

var router *routeManager

func init() {
	iphlpapi := syscall.NewLazyDLL("iphlpapi.dll")
	getIpForwardTable := iphlpapi.NewProc("GetIpForwardTable")
	createIpForwardEntry := iphlpapi.NewProc("CreateIpForwardEntry")
	deleteIpForwardEntry := iphlpapi.NewProc("DeleteIpForwardEntry")

	router = &routeManager{
		iphlpapi:             iphlpapi,
		getIpForwardTable:    getIpForwardTable,
		createIpForwardEntry: createIpForwardEntry,
		deleteIpForwardEntry: deleteIpForwardEntry,
		netInterfaces:        make(map[int]netInterface),
	}
}

func List() ([]Route, error) {
	return router.ListRoutes()
}

func Add(route Route) error {
	rules, err := router.listRoutes()
	if err != nil {
		return err
	}

	var aim *routeRule
	for rule := range rules {
		r := rules[rule]
		if r.Type != 3 {
			continue
		}

		if router.getInterfaceIP(int(r.IfIndex)).ip.String() != route.InterfaceIP {
			continue
		}
		aim = &r
		break
	}

	if aim == nil {
		return fmt.Errorf("don not find interface %s", route.InterfaceIP)
	}

	aim.Metric1 = routeMetric
	aim.Type = 4
	aim.Proto = 3

	if n := copy(aim.Dest[:], []byte(net.ParseIP(route.Destination).To4())); n != 4 {
		return fmt.Errorf("copy destination failed %v != 4", n)
	}
	mask := net.ParseIP(route.Mask).To4()
	if n := copy(aim.Mask[:], []byte(mask)); n != 4 {
		return fmt.Errorf("copy mask %v failed %v != 4", mask, n)
	}
	return router.addRoute(aim)
}

func Delete(rr Route) error {
	rs, err := router.listRoutes()
	if err != nil {
		return err
	}

	var aim *routeRule
	for i := range rs {
		r := rs[i]
		if r.Type != 3 {
			continue
		}

		if router.getInterfaceIP(int(r.IfIndex)).ip.String() != rr.InterfaceIP {
			continue
		}

		if net.IPv4(r.Dest[0], r.Dest[1], r.Dest[2], r.Dest[3]).String() != rr.Destination {
			continue
		}

		if net.IPv4(r.Mask[0], r.Mask[1], r.Mask[2], r.Mask[3]).String() != rr.Mask {
			continue
		}
		aim = &r
		break
	}

	if aim == nil {
		return errors.New("don not find route")
	}
	return router.delRoute(aim)
}

func (rt *routeManager) listRoutes() ([]routeRule, error) {
	buf := make([]byte, 4+unsafe.Sizeof(routeRule{}))
	buf_len := uint32(len(buf))

	rt.getIpForwardTable.Call(uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&buf_len)), 0)

	var r1 uintptr
	for i := 0; i < 5; i++ {
		buf = make([]byte, buf_len)
		r1, _, _ = rt.getIpForwardTable.Call(uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&buf_len)), 0)
		if r1 == 122 {
			continue
		}
		break
	}

	if r1 != 0 {
		return nil, fmt.Errorf("failed to get the routing table, return value: %v", r1)
	}

	num := *(*uint32)(unsafe.Pointer(&buf[0]))
	routes := make([]routeRule, num)
	sr := uintptr(unsafe.Pointer(&buf[0])) + unsafe.Sizeof(num)
	rowSize := unsafe.Sizeof(routeRule{})

	// 安全检查
	if len(buf) < int((unsafe.Sizeof(num) + rowSize*uintptr(num))) {
		return nil, fmt.Errorf("system error: GetIpForwardTable returns the number is too long, beyond the buffer.")
	}

	for i := uint32(0); i < num; i++ {
		routes[i] = *((*routeRule)(unsafe.Pointer(sr + (rowSize * uintptr(i)))))
	}

	return routes, nil
}

func (rt *routeManager) delRoute(rr *routeRule) error {
	r1, _, err := rt.deleteIpForwardEntry.Call(uintptr(unsafe.Pointer(rr)))
	if r1 != 0 {
		return fmt.Errorf("delete routing table%#v error, return value: %v ,err:%v", rr, r1, err)
	}
	return nil
}

type netInterface struct {
	ip   net.IP
	name string
}

func (rt *routeManager) getInterfaceIP(index int) *netInterface {
	if ip, ok := rt.netInterfaces[index]; ok {
		return &ip
	}

	netInf, err := net.InterfaceByIndex(index)
	if err != nil {
		return nil
	}
	addrs, err := netInf.Addrs()
	if err != nil {
		return nil
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip.IsGlobalUnicast() && (ip.To4() != nil || ip.To16() != nil) {
			ret := &netInterface{ip: ip, name: netInf.Name}
			rt.netInterfaces[index] = *ret
			return ret
		}
	}
	return nil

}

func (rt *routeManager) ListRoutes() ([]Route, error) {
	buf := make([]byte, 4+unsafe.Sizeof(routeRule{}))
	buf_len := uint32(len(buf))

	rt.getIpForwardTable.Call(uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&buf_len)), 0)

	var r1 uintptr
	for i := 0; i < 5; i++ {
		buf = make([]byte, buf_len)
		r1, _, _ = rt.getIpForwardTable.Call(uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&buf_len)), 0)
		if r1 == 122 {
			continue
		}
		break
	}

	if r1 != 0 {
		return nil, fmt.Errorf("failed to get the routing table, return value: %v", r1)
	}

	num := *(*uint32)(unsafe.Pointer(&buf[0]))
	routes := make([]routeRule, num)
	sr := uintptr(unsafe.Pointer(&buf[0])) + unsafe.Sizeof(num)
	rowSize := unsafe.Sizeof(routeRule{})

	// 安全检查
	if len(buf) < int((unsafe.Sizeof(num) + rowSize*uintptr(num))) {
		return nil, fmt.Errorf("system error: GetIpForwardTable returns the number is too long, beyond the buffer")
	}

	for i := uint32(0); i < num; i++ {
		routes[i] = *((*routeRule)(unsafe.Pointer(sr + (rowSize * uintptr(i)))))
	}

	ret := make([]Route, 0)
	for _, route := range routes {
		netIf := rt.getInterfaceIP(int(route.IfIndex))
		ifIp := ""
		ifName := ""
		if netIf != nil {
			ifName = netIf.name
			if netIf.ip != nil {
				ifIp = netIf.ip.String()
			}
		}
		ret = append(ret, Route{
			Destination:   net.IPv4(route.Dest[0], route.Dest[1], route.Dest[2], route.Dest[3]).String(),
			Mask:          net.IPv4(route.Mask[0], route.Mask[1], route.Mask[2], route.Mask[3]).String(),
			Gateway:       net.IPv4(route.NextHop[0], route.NextHop[1], route.NextHop[2], route.NextHop[3]).String(),
			InterfaceIP:   ifIp,
			InterfaceName: ifName,
		})
	}

	return ret, nil
}

func (rt *routeManager) addRoute(rr *routeRule) error {
	r1, _, err := rt.createIpForwardEntry.Call(uintptr(unsafe.Pointer(rr)))
	fmt.Printf("r:%v,err:%v", r1, err)
	if r1 == 5010 {
		return nil
	} else if r1 != 0 {
		return fmt.Errorf("add routing table%#v error, return value: %v, err:%v", rr, r1, err)
	}
	return nil
}
