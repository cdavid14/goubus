package goubus

import (
	"net"
)

type UbusNetwork struct {
	NetworkInterface NetworkInterface  `json:"network"`
	NetworkDevice    UbusNetworkDevice `json:"network.device"`
}

func (u UbusNetwork) GetDevicesType() string {
	if u.NetworkInterface.Device == "bridge" {
		return "Bridge"
	} else if u.NetworkInterface.Device == "bond" {
		return "Bond"
	} else {
		return "Ethernet"
	}
}

func (u UbusNetwork) GetDeviceStatus() string {
	if u.NetworkDevice.Up {
		return "up"
	} else {
		return "down"
	}
}

func (u UbusNetwork) GetDeviceMethod() string {
	if u.NetworkInterface.Proto == "static" {
		return "manual"
	} else {
		return "auto"
	}
}

type Addresses struct {
	Address string `json:"Address"`
	Prefix  uint8  `json:"Prefix"`
}

type Routes struct {
	Destination          string            `json:"Destination"`
	Prefix               uint8             `json:"Prefix"`
	NextHop              string            `json:"NextHop"`
	Metric               uint8             `json:"Metric"`
	AdditionalAttributes map[string]string `json:"AdditionalAttributes"`
}

type IP4Config struct {
	Addresses   []Addresses `json:"Addresses"`
	Gateway     string      `json:"Gateway"`
	Domains     []string    `json:"Domains"`
	Nameservers []string    `json:"Nameservers"`
	Routes      []Routes    `json:"Routes"`
}

func (u *UbusNetwork) GetDeviceIp() IP4Config {
	ip4Config := IP4Config{}
	for _, ip := range u.NetworkInterface.Ipv4Address {
		address := Addresses{
			Address: ip.Address,
			Prefix:  uint8(ip.Mask),
		}
		ip4Config.Addresses = append(ip4Config.Addresses, address)
	}
	ip4Config.Gateway = u.NetworkInterface.Route[0].Nexthop
	ip4Config.Domains = u.NetworkInterface.DNSServer
	ip4Config.Nameservers = u.NetworkInterface.DNSServer
	for _, route := range u.NetworkInterface.Route {
		r := Routes{
			Destination: route.Target,
			Prefix:      uint8(route.Mask),
			NextHop:     route.Nexthop,
			Metric:      100,
		}
		ip4Config.Routes = append(ip4Config.Routes, r)
	}
	return ip4Config
}

func (u *Ubus) NetworkStatus(id int, infName string) UbusNetwork {
	ubusNetwork := UbusNetwork{}
	networkInterface, _ := u.NetworkInterfaceStatus(id, infName)
	networkDevice, _ := u.NetworkDeviceStatus(id, networkInterface.L3Device)
	ubusNetwork.NetworkInterface = networkInterface
	ubusNetwork.NetworkDevice = networkDevice
	return ubusNetwork
}

func (u *Ubus) NetworkDevices() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var interfaceNames []string
	for _, i := range interfaces {
		interfaceNames = append(interfaceNames, i.Name)
	}
	return interfaceNames, nil
}
