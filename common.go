package goubus

import "time"

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
	AuthData UbusAuthData
}

//LoginExpired check if login RPC Session id has expired
func (u *Ubus) LoginExpired() bool {
	return u.AuthData.ExpireTime.Before(time.Now())
}
