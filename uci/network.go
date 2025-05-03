package uci

import (
	"github.com/digineo/go-uci/v2"
)

type NetworkInterface struct {
	Name    string `json:"name"`
	Type    string `json:"type"`  // "bridge" or "interface"
	Proto   string `json:"proto"` // "static" or "dhcp"
	Ipaddr  string `json:"ipaddr"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Dns     string `json:"dns"`
	Macaddr string
	Device  []string `json:"ifname"`
}

type Network struct {
	Network map[string]NetworkInterface
	UciTree uci.Tree
}

func NewNetworkConfig(ut uci.Tree) *Network {
	if ut == nil {
		ut = uci.NewTree("/etc/config")
	}
	return &Network{
		Network: make(map[string]NetworkInterface),
		UciTree: ut,
	}
}

func (n *Network) AddInterface(name string, iface NetworkInterface) {
	n.Network[name] = iface
}

func (n *Network) GetInterface(name string) NetworkInterface {
	return n.Network[name]
}

func (n *Network) DelInterface(name string) {
	delete(n.Network, name)
}

func (n *Network) GetInterfaces() map[string]NetworkInterface {
	return n.Network
}

func (n *Network) SetInterfaces(ifaces map[string]NetworkInterface) {
	n.Network = ifaces
}

func (n *Network) GetInterfaceNames() []string {
	names := []string{}
	for name := range n.Network {
		names = append(names, name)
	}
	return names
}

func (n *Network) SaveConfig() error {
	for _, iface := range n.Network {
		err := n.ReplayInterface(iface)
		if err != nil {
			return err
		}
	}

	return n.UciTree.Commit()
}

func (n *Network) ReplayInterface(inf NetworkInterface) error {
	if inf.Type != "bridge" && len(inf.Device) == 0 {
		inf.Device = []string{inf.Name}
	}
	infMap, err := StructToMap(inf)
	if err != nil {
		return err
	}

	err = n.UciTree.DelSection("network", inf.Name)
	if err != nil {
		return err
	}
	err = n.UciTree.AddSection("network", inf.Name, "interface")
	if err != nil {
		return err
	}
	for k, v := range infMap {
		if k == "ifname" {
			if len(inf.Device) > 1 {
				n.UciTree.SetType("network", inf.Name, "ifname", uci.TypeList, inf.Device...)
			} else if len(inf.Device) == 1 {
				n.UciTree.SetType("network", inf.Name, "ifname", uci.TypeOption, inf.Device[0])
			}
			continue
		}

		n.UciTree.SetType("network", inf.Name, k, uci.TypeOption, v)
	}
	n.UciTree.SetType("network", inf.Name, "ipv6", uci.TypeOption, "0")
	return nil
}

func (n *Network) GetInterfacesInConfig() ([]string, error) {
	sections, err := n.UciTree.GetSections("network", "interface")
	if err != nil {
		return nil, err
	}
	return sections, nil
}
