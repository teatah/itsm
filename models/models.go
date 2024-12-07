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

type Conversation struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	User1ID  uint      `json:"user1_id"`
	User2ID  uint      `json:"user2_id"`
	Messages []Message `json:"messages" gorm:"foreignKey:ConversationID"`
}

type Message struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	SenderID       uint   `json:"sender_id"`
	ReceiverID     uint   `json:"receiver_id"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
	ConversationID uint   `json:"conversation_id"`
}
