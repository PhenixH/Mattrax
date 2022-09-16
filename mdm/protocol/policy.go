package protocol

type Policy struct {
	Restrictions *PolicyRestrictions `json:"restrictions"`
	WiFi         *PolicyWiFi         `json:"wifi"`
}

type PolicyRestrictions struct {
	DisableCamera    bool `json:"disable_camera"`
	DisableBluetooth bool `json:"disable_bluetooth"`
}

type PolicyWiFi struct {
	SSID       string `json:"ssid"`
	Passphrase string `json:"passphrase"`
}

// func init() {
// 	var p = Policy{
// 		Restrictions: &PolicyRestrictions{},
// 		WiFi:         &PolicyWiFi{},
// 	}
// 	data, _ := json.Marshal(p)
// 	fmt.Println(string(data))
// }
