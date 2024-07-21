package models

import (
	"time"
)

// Resource struct to describe Resource object.
type Resource struct {
	ID   int `db:"id" json:"id"`
	Name string `db:"name" json:"name" validate:"required"`
	StackId int `db:"project_id" json:"project_id"`
	CloudProvId int  `db:"cloud_prov_id" json:"cloud_prov_id"`
	TypeId int `db:"res_type_id" json:"res_type_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`

}

type CreateResourceRequest struct {
	ResourceName   string `json:"resource_name" validate:"required,min=3,max=32"`
	Type string `json:"type" validate:"required"`
	CloudProviderSlug string `json:"cloud_provider_slug" validate:"required,min=3,max=32"`
	StackName string `json:"stack_name" validate:"required,min=3,max=32"`
}
