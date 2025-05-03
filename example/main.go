package main

import (
	"log"

	"github.com/ploynomail/goubus"
	"github.com/ploynomail/goubus/uci"
)

// This is an example acl of how to use the goubus package,
// you can use the root user to access all the ubus functions,
// and you can increase the access control list to allow other users to access the ubus functions.
// /usr/share/rpcd/acl.d/root.json
// {
//         "root": {
//                 "description": "Super user access role",
//                 "read": {
//                         "ubus": {
//                                 "*": [ "*" ]
//                         },
//                         "uci": [ "*" ],
//                         "file": {
//                                 "*": ["*"]
//                         }
//                 },
//                 "write": {
//                         "ubus": {
//                                 "*": [ "*" ]
//                         },
//                         "uci": [ "*" ],
//                         "file": {
//                                 "*": ["*"]
//                         },
//                         "cgi-io": ["*"]
//                 }
//         }
// }

func main() {
	ubus := goubus.Ubus{
		Username: "root",
		Password: "",
		URL:      "http://192.168.23.196/ubus",
		UciTree:  goubus.NewUciTree("/etc/config"),
	}
	_, err := ubus.AuthLogin()
	if err != nil {
		log.Fatal(err)
	}
	if err := ubus.LoginCheck(); err != nil {
		log.Fatal(err)
	}
	n := uci.NewNetworkConfig(ubus.UciTree)
	c, err := n.GetInterfacesInConfig()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range c {
		networkInfStatus := ubus.NetworkStatus(0, v)
		log.Printf("%+v\n", networkInfStatus)
	}
	// l, err := ubus.LogRead(5, 10, false, false)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("%+v\n", l)
	// var serverRequest = goubus.ServiceListRequest{
	// 	Name:    "uhttpd",
	// 	Verbose: false,
	// }
	// sl, err := ubus.GetServceList(0, serverRequest)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("%+v\n", sl)
	// err = ubus.RcInit(0, goubus.UbusRcInitRequest{
	// 	Name:   "uwsgi",
	// 	Action: "enable",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// rclist, err := ubus.RcList(0)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("%+v\n", rclist)
	// network := uci.Network{
	// 	Network: map[string]uci.NetworkInterface{
	// 		"lan": {
	// 			Name:    "lan",
	// 			Proto:   "static",
	// 			Type:    "bridge",
	// 			Ipaddr:  "192.168.23.196",
	// 			Netmask: "255.255.254.0",
	// 			Dns:     "223.5.5.5",
	// 			Gateway: "192.168.22.1",
	// 			Device:  []string{"eth0", "eth1"},
	// 		},
	// 		"wan": {
	// 			Name:    "wan",
	// 			Proto:   "static",
	// 			Ipaddr:  "192.168.57.1",
	// 			Netmask: "255.255.255.0",
	// 			Gateway: "192.168.57.254",
	// 			Dns:     "192.168.57.254",
	// 			Type:    "ethernet",
	// 			Device:  []string{"eth2"},
	// 		},
	// 	},
	// }
	// n := uci.NewNetworkConfig(ubus.UciTree)
	// n.SetInterfaces(network.Network)
	// err := n.SaveConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c, err := n.GetInterfacesInConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("%+v\n", c)
}
