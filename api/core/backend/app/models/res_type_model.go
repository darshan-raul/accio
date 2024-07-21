package models

// ResourceType struct to describe Resource Type object.
type ResourceType struct {
	ID   int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}