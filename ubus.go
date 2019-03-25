package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//UbusResponseCode represent the status code from JSON-RPC Call
type UbusResponseCode float64

//Represents enum ubus_msg_status from https://git.openwrt.org/?p=project/ubus.git;a=blob;f=ubusmsg.h;h=398b126b6dc01833937749a110181ea0debb1476;hb=HEAD
const (
	UbusStatusOK               UbusResponseCode = 0
	UbusStatusInvalidCommand   UbusResponseCode = 1
	UbusStatusInvalidArgument  UbusResponseCode = 2
	UbusStatusMethodNotFound   UbusResponseCode = 3
	UbusStatusNotFound         UbusResponseCode = 4
	UbusStatusNoData           UbusResponseCode = 5
	UbusStatusPermissionDenied UbusResponseCode = 6
	UbusStatusTimeout          UbusResponseCode = 7
	UbusStatusNotSupported     UbusResponseCode = 8
	UbusStatusUnknownError     UbusResponseCode = 9
	UbusStatusConnectionFailed UbusResponseCode = 10
	UbusStatusLast             UbusResponseCode = 11
)

//Ubus represents information to JSON-RPC Interaction with router
type Ubus struct {
	Username string
	Password string
	URL      string
}

//UbusResponse represents the JSON-RPC Response from router
type UbusResponse struct {
	JSONRPC          string
	ID               int
	Result           interface{}
	UbusResponseCode UbusResponseCode
	UbusData         UbusData
}

//UbusData represents the Data response from JSON-RPC Call
type UbusData struct {
	UbusRPCSession string `json:"ubus_rpc_session"`
	Timeout        int
	Expires        int
	ExpireTime     time.Time
	ACLs           UbusACLS `json:"acls"`
	Data           map[string]string
}

//UbusACLS represents the ACL from user on Authentication
type UbusACLS struct {
	AccessGroup map[string][]string `json:"access-group"`
	Ubus        map[string][]string
	Uci         map[string][]string
}

//Login Call JSON-RPC method to Router Authentication
func (u *Ubus) Login() (UbusResponse, error) {
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": 1, 
			"method": "call", 
			"params": [ 
				"00000000000000000000000000000000", 
				"session", 
				"login", 
				{ 
					"username": "` + u.Username + `", 
					"password": "` + u.Password + `"  
				} 
			] 
		}`)
	req, err := http.NewRequest("POST", u.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return UbusResponse{}, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	result := UbusResponse{}
	json.Unmarshal([]byte(body), &result)
	fmt.Println(string(body))
	//Workaround cause response code not contempled by unmarshal function
	result.UbusResponseCode = UbusResponseCode(result.Result.([]interface{})[0].(float64))
	//Workaround to get UbusData cause the structure of this array has a problem with unmarshal
	if result.UbusResponseCode == UbusStatusOK {
		ubusData := UbusData{}
		ubusDataByte, _ := json.Marshal(result.Result.([]interface{})[1])
		json.Unmarshal(ubusDataByte, &ubusData)
		result.UbusData = ubusData
	} else {
		return UbusResponse{}, fmt.Errorf("Ubus Status Failed: %v", result.UbusResponseCode)
	}
	//Editing expire to current time + expire time to check in functions if the session has expired
	result.UbusData.ExpireTime = time.Now().Add(time.Second * time.Duration(result.UbusData.Expires))
	//Clear result Interface data
	result.Result = ""
	return result, nil
}
