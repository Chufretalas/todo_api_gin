package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Passhash string
	TODOs    []TODO `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags     []Tag  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TODO struct {
	gorm.Model
	Title       string
	Description string
	Done        bool
	UserID      uint
	Tags        []Tag `gorm:"many2many:todo_tags;"`
}

type Tag struct {
	gorm.Model
	Name   string
	UserID uint
}
