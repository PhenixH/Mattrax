package main

import (
	"strings"
	"testing"
)

var testUPNs = map[string]string{
	"not-an-email":               "",
	"test@email@otbeaumont.me":   "otbeaumont.me",
	"test@otbeaumont.me":         "otbeaumont.me",
	"test@student.otbeaumont.me": "student.otbeaumont.me",
}

func TestUPNParsing(t *testing.T) {
	for upn, propperDomain := range testUPNs {
		domain, err := parseUPN(upn)
		if propperDomain == "" {
			if err == nil || !strings.HasPrefix(err.Error(), "error invalid upn format") {
				t.Errorf("Invalid email was parsed correctly, email: %v", domain)
			}
		} else {
			if err == nil && domain != propperDomain {
				t.Errorf("Email domain parsed incorrectly, domain: %v, want: %v.", domain, propperDomain)
			}
		}
	}
}
