// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// AccountResource provides business logic for account operations.
type AccountResource struct {
	BaseResource
	service services.AccountServicer
}

// NewAccountResource creates a new AccountResource with the given service.
func NewAccountResource(svc services.AccountServicer) AccountResourcer {
	return &AccountResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetAll retrieves all accounts from the API.
// This is a pass-through to the service layer for pure API access.
func (r *AccountResource) GetAll() ([]services.Account, error) {
	return r.service.GetAll()
}

// Get retrieves a specific account by ID from the API.
// This is a pass-through to the service layer for pure API access.
func (r *AccountResource) Get(id string) (*services.Account, error) {
	return r.service.Get(id)
}

// Activate activates an account by ID.
// This is a pass-through to the service layer for pure API access.
func (r *AccountResource) Activate(id string) error {
	return r.service.Activate(id)
}

// Deactivate deactivates an account by ID.
// This is a pass-through to the service layer for pure API access.
func (r *AccountResource) Deactivate(id string) error {
	return r.service.Deactivate(id)
}

// GetByName retrieves an account by username using client-side filtering.
// It fetches all accounts and searches for a matching username.
func (r *AccountResource) GetByName(name string) (*services.Account, error) {
	logging.Trace()

	accounts, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return FindByName(accounts, "account", name, func(a services.Account) string {
		return a.Username
	})
}
