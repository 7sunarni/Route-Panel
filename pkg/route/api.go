package route

type Route struct {
	Destination   string `json:"destination"`
	Mask          string `json:"mask"`
	Gateway       string `json:"gateway"`
	InterfaceIP   string `json:"interfaceIp"`
	InterfaceName string `json:"interfaceName"`
}

func (r *Route) IsDefaultRoute() bool {
	return r.Destination == "0.0.0.0" &&
		r.Mask == "0.0.0.0"
}

type RouteError string

const (
	InterfaceNotExist   RouteError = "do not find interface"
	CopyAddressFailed   RouteError = "copy ipv4 address"
	RouteNotExist       RouteError = "do not find route"
	GetRouteTableFailed RouteError = "get routing table failed"
	SyscallFailed       RouteError = "syscall size too big"
	DeleteFailed        RouteError = "delete route failed"
)
