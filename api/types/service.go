package types

// Port estructura que encapsula la información relacionada a un puerto de un contenedor
type Port struct {
	Publics []string `json:"publics"`
	Type    string   `json:"protocol"`
}

type Container struct {
	ID     string          `json:"id"`
	HostIP string          `json:"host_ip,omitempty"`
	Ports  map[string]Port `json:"ports,omitempty"`
	Status string          `json:"status,omitempty"`
	Name   string          `json:"name,omitempty"`
	Image  string          `json:"image,omitempty"`
}
