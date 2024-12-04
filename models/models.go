package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:"default:false"`
	IsTech   bool   `gorm:"default:false"`
}

type ServiceLine struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type Service struct {
	ID            uint        `gorm:"primaryKey"`
	ServiceLineID uint        `gorm:"not null"`
	Name          string      `gorm:"not null"`
	Description   string      `json:"description"`
	ServiceLine   ServiceLine `gorm:"foreignKey:ServiceLineID"`
}

type Incident struct {
	ID      uint   `gorm:"primaryKey"`
	Details string `gorm:"not null"`
}

type Message struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"not null"`
}
