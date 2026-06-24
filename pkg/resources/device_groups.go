// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// DeviceGroupResource provides business logic for device group operations.
type DeviceGroupResource struct {
	BaseResource
	service services.DeviceGroupServicer
}

// NewDeviceGroupResource creates a new DeviceGroupResource with the given service.
func NewDeviceGroupResource(svc services.DeviceGroupServicer) DeviceGroupResourcer {
	return &DeviceGroupResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetByName retrieves a device group by name using client-side filtering.
// It fetches all device groups and searches for a matching name.
func (r *DeviceGroupResource) GetByName(name string) (*services.DeviceGroup, error) {
	logging.Trace()

	groups, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return FindByName(groups, "device group", name, func(g services.DeviceGroup) string {
		return g.Name
	})
}
