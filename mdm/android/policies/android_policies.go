package policies

import (
	"encoding/json"
	"fmt"

	"github.com/mattrax/Mattrax/mdm/protocol"
	"google.golang.org/api/androidmanagement/v1"
	"google.golang.org/api/googleapi"
)

type OpenNetworkConfiguration struct {
	NetworkConfigurations []OpenNetworkConfig
}

type OpenNetworkConfig struct {
	GUID string
	Name string
	Type string
	WiFi *OpenNetworkConfigWifi
}

type OpenNetworkConfigWifi struct {
	SSID       string
	Security   string
	Passphrase string
}

func GenerateAndroidPolicy(p protocol.Policy) (androidmanagement.Policy, error) {
	var ap = androidmanagement.Policy{}
	if p.Restrictions != nil {
		ap.CameraDisabled = p.Restrictions.DisableCamera
		ap.BluetoothDisabled = p.Restrictions.DisableBluetooth
	}
	if p.WiFi != nil {
		var netcfg = OpenNetworkConfiguration{
			NetworkConfigurations: []OpenNetworkConfig{
				{
					GUID: "TODO",
					Name: "Wifi TODO",
					Type: "WiFi",
					WiFi: &OpenNetworkConfigWifi{
						SSID:       p.WiFi.SSID,
						Security:   "WPA-PSK",
						Passphrase: p.WiFi.Passphrase,
					},
				},
			},
		}

		networkConfigJSON, err := json.Marshal(netcfg)
		if err != nil {
			return ap, fmt.Errorf("error marshalling open network configuration: %w", err)
		}
		ap.OpenNetworkConfiguration = googleapi.RawMessage(networkConfigJSON)
	}

	return ap, nil
}
