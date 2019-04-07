# GoUbus

That library was developed for communication with [Ubus (OpenWrt micro bus architecture)](https://git.openwrt.org/project/ubus.git) with http requests. The library is in **alpha-development** and **NOT RECOMMENDED FOR PRODUCTION** until the moment.

The project was created to solve a problem at UNIVERSIDADE DO ESTADO DO RIO GRANDE DO NORTE - UERN with network management with low-cost routers.

Please check our [License](LICENSE)

# TO-DO

- [X] Structs to represent the Ubus interaction following [JSON-RPC 2.0](https://www.jsonrpc.org/specification)
- [X] Login Function
- [X] Wireless Info Get
- [ ] Interfaces Info Get
- [X] DHCP Info Get
- [X] File Management
- [X] Log System
- [ ] UCI
- [X] Check if the necessary plugins are installed on OpenWRT Router

# Conventions

The ID 1 from JSON-RPC 2.0 Requests is reserved for Login procedure


# Example


    package main
    
    import (
    	"fmt"
    	"log"
    	"github.com/cdavid14/goubus"
    )
    
    func main() {
    	ubus := goubus.Ubus{
    	  Username: "root",
    	  Password: "admin",
    	  URL:      "http://192.168.1.1/ubus",
    	}
    	result, err := ubus.Login()
    	if err != nil {
    	  log.Fatal(err)
    	}
    	fmt.Println(result)
    }

and it will return something like

> {8b93715dbb85378d87daf0b1cc64a83b 300 299 2019-03-25
> 01:36:28.248837063 -0300 -03 m=+299.016411013 {map[uci-access:[read
> write] unauthenticated:[read]] map[session:[access login]] map[*:[read
> write]]} map[username:root]}

# Final Notes

Please contribute to make this library most usually as possible and improve more functions!

