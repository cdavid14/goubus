package goubus

import (
	"encoding/json"
	"errors"
	"strconv"
)

type NetworkInterface struct {
	Up          bool     `json:"up"`
	Pending     bool     `json:"pending"`
	Available   bool     `json:"available"`
	Autostart   bool     `json:"autostart"`
	Dynamic     bool     `json:"dynamic"`
	Uptime      int      `json:"uptime"`
	L3Device    string   `json:"l3_device"`
	Proto       string   `json:"proto"`
	Device      string   `json:"device"`
	Updated     []string `json:"updated"`
	Metric      int      `json:"metric"`
	DNSMetric   int      `json:"dns_metric"`
	Delegation  bool     `json:"delegation"`
	Ipv4Address []struct {
		Address string `json:"address"`
		Mask    int    `json:"mask"`
	} `json:"ipv4-address"`
	Ipv6Address []struct {
		Address string `json:"address"`
		Mask    int    `json:"mask"`
	} `json:"ipv6-address"`
	Ipv6Prefix           []interface{} `json:"ipv6-prefix"`
	Ipv6PrefixAssignment []interface{} `json:"ipv6-prefix-assignment"`
	Route                []struct {
		Target  string `json:"target"`
		Mask    int    `json:"mask"`
		Nexthop string `json:"nexthop"`
		Source  string `json:"source"`
	} `json:"route"`
	DNSServer []string      `json:"dns-server"`
	DNSSearch []interface{} `json:"dns-search"`
	Neighbors []interface{} `json:"neighbors"`
	Inactive  struct {
		Ipv4Address []interface{} `json:"ipv4-address"`
		Ipv6Address []interface{} `json:"ipv6-address"`
		Route       []interface{} `json:"route"`
		DNSServer   []interface{} `json:"dns-server"`
		DNSSearch   []interface{} `json:"dns-search"`
		Neighbors   []interface{} `json:"neighbors"`
	} `json:"inactive"`
	Data struct {
	} `json:"data"`
}

func (u *Ubus) NetworkInterfaceStatus(id int, name string) (NetworkInterface, error) {
	errLogin := u.LoginCheck()
	if errLogin != nil {
		return NetworkInterface{}, errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.AuthData.UbusRPCSession + `", 
				"network.interface.` + name + `",
				"status", 
				{} 
			] 
		}`)
	call, err := u.Call(jsonStr)
	if err != nil {
		return NetworkInterface{}, err
	}
	ubusData := NetworkInterface{}

	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return NetworkInterface{}, errors.New("data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}
