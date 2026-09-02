// Package state defines types and functions for managing the lifecycle, configuration, and resource tracking of sandboxes.
package state

import (
	"errors"
	"slices"
	"time"
)

// ErrInvalidState is returned when a Desired or Observed state is invalid.
var ErrInvalidState = errors.New("state: desired or observed state is not a valid state")

type DesiredState string

const (
	DesiredRunning DesiredState = "running"
	DesiredStopped DesiredState = "stopped"
	DesiredDeleted DesiredState = "deleted"
)

type ObservedState string

const (
	ObservedRequested    ObservedState = "requested"
	ObservedProvisioning ObservedState = "provisioning"
	ObservedPrepared     ObservedState = "prepared"
	ObservedStarting     ObservedState = "starting"
	ObservedRunning      ObservedState = "running"
	ObservedStopping     ObservedState = "stopping"
	ObservedStopped      ObservedState = "stopped"
	ObservedDeleting     ObservedState = "deleting"
	ObservedFailed       ObservedState = "failed"
)

// Sandbox describes an isolated sandbox and its allocated resources.
type Sandbox struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Slot          int           `json:"slot"`
	UID           int           `json:"uid"`
	GID           int           `json:"gid"`
	VCPUs         int           `json:"vcpus"`
	MemoryMiB     int           `json:"memory_mib"`
	GuestIP       string        `json:"guest_ip"`
	GatewayIP     string        `json:"gateway_ip"`
	GuestSubnet   string        `json:"guest_subnet"`
	TransitHostIP string        `json:"transit_host_ip"`
	TransitNSIP   string        `json:"transit_ns_ip"`
	TransitSubnet string        `json:"transit_subnet"`
	Namespace     string        `json:"namespace"`
	HostVeth      string        `json:"host_veth"`
	NSVeth        string        `json:"ns_veth"`
	TAP           string        `json:"tap"`
	DesiredState  DesiredState  `json:"desired_state"`
	ObservedState ObservedState `json:"observed_state"`
	Generation    uint64        `json:"generation"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	LastError     string        `json:"last_error"`
}

func (d DesiredState) Valid() bool {
	desiredStates := []DesiredState{DesiredRunning, DesiredStopped, DesiredDeleted}

	return slices.Contains(desiredStates, d)
}

func (o ObservedState) Valid() bool {
	observedStates := []ObservedState{
		ObservedRequested,
		ObservedProvisioning,
		ObservedPrepared,
		ObservedStarting,
		ObservedRunning,
		ObservedStopping,
		ObservedStopped,
		ObservedDeleting,
		ObservedFailed,
	}

	return slices.Contains(observedStates, o)
}
