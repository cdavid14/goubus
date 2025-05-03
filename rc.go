package goubus

import (
	"encoding/json"
	"errors"
	"strconv"
)

type UbusRcInitRequest struct {
	Name   string `json:"name"`
	Action string `json:"action"`
}

type UbusRcListResponse struct {
	Start   int  `json:"start"`
	Running bool `json:"running"`
	Enabled bool `json:"enabled"`
}

func (u *Ubus) RcList(id int) (map[string]UbusRcListResponse, error) {
	var resp map[string]UbusRcListResponse = make(map[string]UbusRcListResponse)
	errLogin := u.LoginCheck()

	if errLogin != nil {
		return resp, errLogin
	}
	var err error
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.AuthData.UbusRPCSession + `", 
				"rc", 
				"list", 
				{}
			] 
		}`)
	call, err := u.Call(jsonStr)
	if err != nil {
		return resp, err
	}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return resp, errors.New("data error")
	}
	if err := json.Unmarshal(ubusDataByte, &resp); err != nil {
		return resp, errors.New("data error")
	}
	return resp, nil
}

func (u *Ubus) RcInit(id int, request UbusRcInitRequest) error {
	errLogin := u.LoginCheck()

	if errLogin != nil {
		return errLogin
	}
	var jsonData []byte
	var err error
	if request.Name == "" || request.Action == "" {
		return errors.New("name and action required")
	}
	jsonData, err = json.Marshal(request)
	if err != nil {
		return errors.New("error parsing UCI request data")
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.AuthData.UbusRPCSession + `", 
				"rc", 
				"init", 
				` + string(jsonData) + `
			] 
		}`)
	call, err := u.Call(jsonStr)
	if err != nil {
		return err
	}
	if call.Error.Code != 0 {
		return errors.New(call.Error.Message)
	}
	return nil
}
