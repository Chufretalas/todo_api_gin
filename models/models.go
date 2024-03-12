package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"column:username"`
	Passhash string `gorm:"column:pashash"`
	TODOs    []TODO
	Tags     []Tag
}

type TODO struct {
	gorm.Model
	Title       string
	Description string
	UserID      uint
	// User        User  `gorm:"foreignKey:UserID,constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tags []Tag `gorm:"many2many:todo_tags;"`
}

type Tag struct {
	gorm.Model
	Name   string
	UserID uint
	// User   User `gorm:"foreignKey:UserID,constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
