package goubus

import (
	"log"
	"testing"
)

func TestNetworkDevices(t *testing.T) {

	u := &Ubus{}
	devices, err := u.NetworkDevices()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	log.Printf("%+v\n", devices)
}
