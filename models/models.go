package models

type User struct {
	ID               uint   `gorm:"primaryKey"`
	Username         string `gorm:"not null"`
	Password         string `gorm:"not null"`
	IsAdmin          bool   `gorm:"default:false"`
	IsTechOfficer    bool   `gorm:"default:false"`
	IsDefaultOfficer bool   `gorm:"default:false"`
}

type Service struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `json:"description"`
	IsBusiness  bool   `gorm:"default:false"`
	IsTechnical bool   `gorm:"default:false"`
}

type Incident struct {
	ID      uint   `gorm:"primaryKey"`
	Details string `gorm:"not null"`
}

type Message struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"not null"`
}
