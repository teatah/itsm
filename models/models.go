package models

import "time"

type User struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	Username         string `gorm:"not null" json:"username"`
	Password         string `gorm:"not null" json:"password"`
	IsAdmin          bool   `gorm:"default:false" json:"is_admin"`
	IsTechOfficer    bool   `gorm:"default:false" json:"is_tech_officer"`
	IsDefaultOfficer bool   `gorm:"default:false" json:"is_default_officer"`
}

type Service struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `json:"description"`
	IsBusiness  bool   `gorm:"default:false"`
	IsTechnical bool   `gorm:"default:false"`
}

type Incident struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Status      string    // Например, "open", "in progress", "resolved"
	UserID      uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoCreateTime"`
}

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DialogID   uint      `json:"dialog_id" gorm:"not null;foreignKey:DialogID"`
	SenderID   uint      `json:"sender_id" gorm:"not null;foreignKey:UserID"`
	ReceiverID uint      `json:"receiver_id" gorm:"not null;foreignKey:UserID"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	Timestamp  time.Time `json:"timestamp" gorm:"timestamp"`
}

type Dialog struct {
	ID            uint `json:"id"`
	User1ID       uint `json:"user1_id" gorm:"foreignKey:UserID"`
	User2ID       uint `json:"user2_id" gorm:"foreignKey:UserID"`
	LastMessageID uint `json:"last_message_id" gorm:"foreignKey:MessageID"`
}
