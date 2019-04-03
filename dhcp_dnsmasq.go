package goubus

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type DHCPIPv4Lease struct {
	Hostname  string
	IPv4      string
	MAC       string
	Leasetime time.Time
	ClientID  string
}

type DHCPIPv4Leases []DHCPIPv4Lease

func (u *Ubus) DHCPIPv4Leases(id int) (DHCPIPv4Leases, error) {
	errLogin := u.LoginCheck()
	if errLogin != nil {
		return DHCPIPv4Leases{}, errors.New("Error on Login")
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.AuthData.UbusRPCSession + `", 
				"file", 
				"exec", 
				{ 
					"command": "cat",
					"params": ['/tmp/dhcp.leases']
				} 
			] 
		}`)
	call, err := u.Call(jsonStr)
	if err != nil {
		return DHCPIPv4Leases{}, err
	}
	DHCPData := DHCPIPv4Leases{}
	ubusData := UbusExec{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return DHCPIPv4Leases{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	for _, line := range strings.Split(ubusData.Stdout, "\n") {
		if len(line) > 5 {
			data := strings.Split(line, " ")
			timedata, _ := strconv.Atoi(data[0])
			DHCPData = append(DHCPData, DHCPIPv4Lease{
				Leasetime: time.Unix(int64(timedata), 0),
				MAC:       data[1],
				IPv4:      data[2],
				Hostname:  data[3],
				ClientID:  data[4],
			})
		}
	}
	return DHCPData, nil
}
