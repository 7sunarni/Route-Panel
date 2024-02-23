package route

type Route struct {
	Destination   string `json:"destination"`
	Mask          string `json:"mask"`
	Gateway       string `json:"gateway"`
	InterfaceIP   string `json:"interfaceIp"`
	InterfaceName string `json:"interfaceName"`
	Metric        string `json:"metric"`
	Type          string `json:"type"`
	Protocol      string `json:"protocol"`
}

func (r *Route) IsDefaultRoute() bool {
	return r.Type == "Indirect" &&
		r.Protocol == "Static Route" &&
		r.Destination == "0.0.0.0" &&
		r.Mask == "0.0.0.0"
}
