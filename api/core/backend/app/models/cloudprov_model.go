package models

// CloudProv struct to describe CloudProv object.
type CloudProv struct {
	ID   int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`
}