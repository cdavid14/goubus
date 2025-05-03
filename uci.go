package goubus

import "github.com/digineo/go-uci/v2"

var defaultUciPath = "/etc/config"

type UbusUciRequestGeneric struct {
	Config  string `json:"config"`
	Section string `json:"section,omitempty"`
	Option  string `json:"option,omitempty"`
	Type    string `json:"type,omitempty"`
	Match   string `json:"match,omitempty"`
	Name    string `json:"name,omitempty"`
}

type UbusUciRequest struct {
	UbusUciRequestGeneric
	Values map[string]string `json:"values,omitempty"`
}

func NewUciTree(defalt string) uci.Tree {
	if defalt != "" {
		defaultUciPath = defalt
	}
	return uci.NewTree(defaultUciPath)
}
