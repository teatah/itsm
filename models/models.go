package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:"default:false"`
}

type BusinessService struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type TechnicalService struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type Incident struct {
	ID      uint   `gorm:"primaryKey"`
	Details string `gorm:"not null"`
}

type Message struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"not null"`
}
