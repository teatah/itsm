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
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	IsBusiness  bool   `gorm:"default:false" json:"is_business"`
	IsTechnical bool   `gorm:"default:false" json:"is_technical"`
}

type Incident struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"not null" json:"user_id"`
	ResponsibleUserID *uint     `gorm:"default:null" json:"responsible_user_id"`
	Title             string    `gorm:"not null" json:"title"`
	Description       string    `json:"description"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	User              User      `gorm:"foreignKey:UserID" json:"user"`
	ResponsibleUser   User      `gorm:"foreignKey:ResponsibleUserID" json:"responsible_user"`
}

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DialogID   uint      `json:"dialog_id" gorm:"not null"`
	SenderID   uint      `json:"sender_id" gorm:"not null"`
	ReceiverID uint      `json:"receiver_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	Timestamp  time.Time `json:"timestamp" gorm:"timestamp"`
	Dialog     Dialog    `gorm:"foreignKey:DialogID" json:"dialog"`
	Sender     User      `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver   User      `gorm:"foreignKey:ReceiverID" json:"receiver"`
}

type Dialog struct {
	ID      uint `json:"id" gorm:"primaryKey"`
	User1ID uint `json:"user1_id" gorm:"not null"`
	User2ID uint `json:"user2_id" gorm:"not null"`
	User1   User `gorm:"foreignKey:User1ID" json:"user1"`
	User2   User `gorm:"foreignKey:User2ID" json:"user2"`
}
