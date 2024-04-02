package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Passhash string
	TODOs    []TODO `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags     []Tag  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

//TODO: this might need a on update cascade on Tags, but I am not sure yet
type TODO struct {
	gorm.Model
	Title       string
	Description string
	Done        bool `gorm:"default:false"`
	UserID      uint
	Tags        []Tag `gorm:"many2many:todo_tags;"`
}

type Tag struct {
	gorm.Model
	Name   string
	UserID uint
}
