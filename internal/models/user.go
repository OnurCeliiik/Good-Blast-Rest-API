package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username string    `json:"username" gorm:"uniques;not null"`
	Coins    int       `json:"coins" gorm:"default:1000"`
	Level    int       `json:"level" gorm:"default:1"`
	Country  string    `json:"country" gorm:"not null;default:'Unkown'"`
}
