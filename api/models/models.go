package models

import (

	"gorm.io/gorm"
)

type Resource struct {
	gorm.Model
	Name  string
	TypeID uint
	Type Type
}


type Type struct {
	gorm.Model
	Name string
	CloudID uint
	Cloud Cloud

}

type Cloud struct {

	gorm.Model
	Name string

}

