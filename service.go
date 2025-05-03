package goubus

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Respawn struct {
	Threshold int `json:"threshold"`
	Timeout   int `json:"timeout"`
	Retries   int `json:"retries"`
}

// AvahiInstance represents an instance of the avahi-daemon service.
type ServiceInstance struct {
	Running     bool     `json:"running"`
	Command     []string `json:"command"`
	TermTimeout int      `json:"term_timeout"`
	ExitCode    int      `json:"exit_code"`
	Respawn     Respawn  `json:"respawn"`
}

type ServiceInstanceList struct {
	Instances []ServiceInstance `json:"instances"`
}

// ServiceStatus represents the status of a service.
type ServiceStatus struct {
	Service ServiceInstanceList `json:"service"`
}

type ServiceListRequest struct {
	Name    string `json:"name"`
	Verbose bool   `json:"verbose"`
}

func (u *Ubus) GetServceList(id int, request ServiceListRequest) (ServiceStatus, error) {
	errLogin := u.LoginCheck()
	if errLogin != nil {
		return ServiceStatus{}, errLogin
	}
	var jsonData []byte
	var err error
	if request.Name == "" {
		jsonData = []byte(`{}`)
	} else {
		jsonData, err = json.Marshal(request)
		if err != nil {
			return ServiceStatus{}, errors.New("Error Parsing UCI Request Data")
		}
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.AuthData.UbusRPCSession + `", 
				"service", 
				"list", 
				` + string(jsonData) + `
			] 
		}`)
	call, err := u.Call(jsonStr)
	if err != nil {
		return ServiceStatus{}, err
	}
	if len(call.Result.([]interface{})) <= 1 {
		return ServiceStatus{}, errors.New("Empty response")
	}
	ubusData := ServiceStatus{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return ServiceStatus{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}
