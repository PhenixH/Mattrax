package android

import "github.com/mattrax/Mattrax/internal/db"

func MattraxDeviceState(deviceState string) (state db.DeviceState) {
	if deviceState == "ACTIVE" {
		state = db.DeviceStateManaged
	} else if deviceState == "DISABLED" {
		state = db.DeviceStateDisabled
	} else if deviceState == "DELETED" {
		state = db.DeviceStateUserUnenrolled
	} else if deviceState == "PROVISIONING" {
		state = db.DeviceStateDeploying
	}
	return state
}

func AndroidDeviceState(deviceState db.DeviceState) (state string) {
	if deviceState == db.DeviceStateManaged {
		state = "ACTIVE"
	} else if deviceState == db.DeviceStateDisabled {
		state = "DISABLED"
	} else if deviceState == db.DeviceStateUserUnenrolled {
		state = "DELETED"
	} else if deviceState == db.DeviceStateDeploying {
		state = "PROVISIONING"
	}
	return state
}

func MattraxDeviceOwnership(deviceOwnership string) (ownership db.DeviceOwnership) {
	if deviceOwnership == "COMPANY_OWNED" {
		ownership = db.DeviceOwnershipCorporate
	} else if deviceOwnership == "PERSONALLY_OWNED" {
		ownership = db.DeviceOwnershipPersonal
	}
	return ownership
}

func MattraxManagementScope(deviceScope string) (scope db.ManagementScope) {
	if deviceScope == "DEVICE_OWNER" {
		scope = db.ManagementScopeDevice
	} else if deviceScope == "PROFILE_OWNER" {
		scope = db.ManagementScopeAfwProfile
	}
	return scope
}
