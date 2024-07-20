package models

// CloudProv struct to describe User object.
type CloudProv struct {
	ID   int `db:"id" json:"id"` // validate:"required"`
	Name string `db:"name" json:"name"` // validate:"required"`
	Slug string `db:"slug" json:"slug"` // validate:"required"`
}