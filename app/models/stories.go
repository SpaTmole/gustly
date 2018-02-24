package models

import (
	"github.com/jinzhu/gorm"
	// "github.com/revel/revel"
)

type Story struct {
	gorm.Model
	Type        string
	Content     string
	Views       []*User `json:"-" gorm:"many2many:views;association_jointable_foreignkey:viewed_user_id`
	UserID      uint    `json:owner`
	NextStoryID *uint   `json:next_story_id`
}
