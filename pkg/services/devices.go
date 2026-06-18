// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
)

// Device represents a device in the Configuration Manager
type Device struct {
	OsType     string                 `json:"ostype"`
	DeviceType string                 `json:"device-type"`
	Name       string                 `json:"name"`
	Host       string                 `json:"host"`
	Actions    []string               `json:"actions"`
	Origins    any                    `json:"origins"`
	Properties map[string]interface{} `json:"properties"`
}

// DeviceService provides methods for managing devices
type DeviceService struct {
	BaseService
}

// NewDeviceService creates a new DeviceService with the given client
func NewDeviceService(c client.Client) *DeviceService {
	return &DeviceService{
		BaseService: NewBaseService(c),
	}
}

// unmarshal converts a map to a Device struct, handling dynamic properties
func (svc *DeviceService) unmarshal(in map[string]interface{}, d *Device) error {
	fields := reflect.TypeOf((*Device)(nil)).Elem()

	var keys = []string{}

	for i := 0; i < fields.NumField(); i++ {
		f := fields.Field(i)
		keys = append(keys, f.Tag.Get("json"))
	}

	exists := func(s string) bool {
		for _, ele := range keys {
			if ele == s {
				return true
			}
		}
		return false
	}

	if err := Unmarshal(in, &d); err != nil {
		return err
	}

	properties := map[string]interface{}{}
	for key, value := range in {
		if !exists(key) {
			properties[key] = value
		}
	}

	d.Properties = properties

	return nil
}

// GetAll retrieves all devices from the Configuration Manager
func (svc *DeviceService) GetAll() ([]Device, error) {
	logging.Trace()

	var devices []Device

	var limit = 100
	var skip = 0
	var start = 0

	type Response struct {
		Entity            string                   `json:"entity"`
		Total             int                      `json:"total"`
		TotalByAdapter    map[string]int           `json:"totalByAdapter"`
		UniqueDeviceCount int                      `json:"unique_device_count"`
		ReturnCount       int                      `json:"return_count"`
		StartIndex        int                      `json:"start_index"`
		List              []map[string]interface{} `json:"list"`
	}

	for {
		// Declare res inside the loop so each page decodes into a freshly
		// allocated struct. Reusing a single res across pages lets
		// encoding/json merge map fields and reuse slice backing arrays,
		// bleeding fields from one page's elements into the next.
		var res Response

		body := map[string]interface{}{
			"options": map[string]interface{}{
				"order": "ascending",
				"sort":  []map[string]interface{}{map[string]interface{}{"name": 1}},
				"start": start,
				"limit": limit,
			},
		}

		if err := svc.PostRequest(&Request{
			uri:                "/configuration_manager/devices",
			body:               &body,
			expectedStatusCode: http.StatusOK,
		}, &res); err != nil {
			return nil, err
		}

		for _, ele := range res.List {
			var d Device
			if err := svc.unmarshal(ele, &d); err != nil {
				return nil, err
			}

			devices = append(devices, d)
		}

		if len(devices) == res.Total {
			break
		}

		skip += limit
	}

	return devices, nil
}

// Get retrieves a device by its name
func (svc *DeviceService) Get(name string) (*Device, error) {
	logging.Trace()

	body := map[string]interface{}{
		"options": map[string]interface{}{
			"order":      "ascending",
			"origins":    true,
			"exactMatch": true,
			"filter":     map[string]interface{}{"name": name},
			"limit":      1,
			"start":      0,
		},
	}

	type Response struct {
		Entity            string                   `json:"entity"`
		Total             int                      `json:"total"`
		TotalByAdapter    map[string]int           `json:"totalByAdapter"`
		UniqueDeviceCount int                      `json:"unique_device_count"`
		ReturnCount       int                      `json:"return_count"`
		StartIndex        int                      `json:"start_index"`
		List              []map[string]interface{} `json:"list"`
	}

	var res Response

	if err := svc.PostRequest(&Request{
		uri:                "/configuration_manager/devices",
		body:               &body,
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	if res.ReturnCount == 0 {
		return nil, errors.New("device not found")
	}

	var device Device

	if err := svc.unmarshal(res.List[0], &device); err != nil {
		return nil, err
	}

	return &device, nil

}
