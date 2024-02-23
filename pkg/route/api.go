package route

type Route struct {
	Destnation string `json:"destination"`
	Mask       string `json:"mask"`
	Gateway    string `json:"gateway"`
	Interface  string `json:"interface"`
	Metric     string `json:"metric"`
	Type       string `json:"type"`
	Protocol   string `json:"protocol"`
}

func (r *Route) IsDefaultRoute() bool {
	return r.Type == "Indirect" &&
		r.Protocol == "Static Route" &&
		r.Destnation == "0.0.0.0" &&
		r.Mask == "0.0.0.0"
}
