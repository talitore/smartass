package models

import (
	"encoding/json"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"time"
)

type Asset struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	Name        string        `json:"name" db:"name"`
	Description nulls.String  `json:"description" db:"description"`
	Url         string        `json:"url" db:"url"`
	Labels      slices.String `json:"labels" db:"labels"`
	UserGuid    string        `json:"user_guid" db:"user_guid"`
}

// String is not required by pop and may be deleted
func (a Asset) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Assets is not required by pop and may be deleted
type Assets []Asset

// String is not required by pop and may be deleted
func (a Assets) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Asset) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Name, Name: "Name"},
		&validators.StringIsPresent{Field: a.Url, Name: "Url"},
		&validators.StringIsPresent{Field: a.UserGuid, Name: "UserGuid"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Asset) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Asset) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
