package goubus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//UbusAuthResponse represents a response from auth module
type UbusAuthResponse struct {
	JSONRPC          string
	ID               int
	Result           interface{}
	UbusResponseCode UbusResponseCode
}

//UbusAuthData represents the Data response from auth module
type UbusAuthData struct {
	UbusRPCSession string `json:"ubus_rpc_session"`
	Timeout        int
	Expires        int
	ExpireTime     time.Time
	ACLs           UbusAuthACLS `json:"acls"`
	Data           map[string]string
}

//UbusAuthACLS represents the ACL from user on Authentication
type UbusAuthACLS struct {
	AccessGroup map[string][]string `json:"access-group"`
	Ubus        map[string][]string
	Uci         map[string][]string
}

//AuthLogin Call JSON-RPC method to Router Authentication
func (u *Ubus) AuthLogin() (UbusAuthResponse, error) {
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
		return UbusAuthResponse{}, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	result := UbusAuthResponse{}
	json.Unmarshal([]byte(body), &result)
	//Workaround cause response code not contempled by unmarshal function
	result.UbusResponseCode = UbusResponseCode(result.Result.([]interface{})[0].(float64))
	//Workaround to get UbusData cause the structure of this array has a problem with unmarshal
	if result.UbusResponseCode == UbusStatusOK {
		ubusData := UbusAuthData{}
		ubusDataByte, err := json.Marshal(result.Result.([]interface{})[1])
		if err != nil {
			return UbusAuthResponse{}, errors.New("Data error")
		}
		json.Unmarshal(ubusDataByte, &ubusData)
		u.AuthData = ubusData
	} else {
		return UbusAuthResponse{}, fmt.Errorf("Ubus Status Failed: %v", result.UbusResponseCode)
	}
	//Editing expire to current time + expire time to check in functions if the session has expired
	u.AuthData.ExpireTime = time.Now().Add(time.Second * time.Duration(u.AuthData.Expires))
	//Clear result Interface data
	result.Result = ""
	//Set AuthData to Ubus struct
	return result, nil
}
