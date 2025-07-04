package models

import (
	"time"

	"github.com/google/uuid"
)

/*
Tournament represents a daily tournament
Each tournament lasts from 00:00 UTC to 23:59 UTC
Users are grouped into 35-player groups
*/

// Represents a tournament event that runs daily.
type Tournament struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	UserCount int       `gorm:"default:0" json:"user_count"`
	MaxUsers  int       `gorm:"default:35" json:"max_users"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string

	//Participants []TournamentParticipant `gorm:"foreignKey:TournamentID" json:"participants"`
}

type TournamentParticipant struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"tour_part_id"`
	TournamentID uuid.UUID `gorm:"type:uuid;not null" json:"tournament_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Level        int       `gorm:"not null;default:0" json:"level"`
}
