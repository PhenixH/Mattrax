package windows

import "github.com/mattrax/Mattrax/internal/db"

func MattraxManagementScope(certStore string) (scope db.ManagementScope) {
	if certStore == "User" {
		scope = db.ManagementScopeUser
	} else if certStore == "Device" {
		scope = db.ManagementScopeDevice
	}
	return scope
}

func MattraxDeviceOwnership(enrollmentType string) (ownership db.DeviceOwnership) {
	if enrollmentType == "Device" {
		ownership = db.DeviceOwnershipCorporate
	} else {
		ownership = db.DeviceOwnershipPersonal
	}
	return ownership
}
