package models

type Backdrop struct {
	BaseModel
	Name     string `gorm:"size:100;not null" json:"name"`
	ImageURL string `gorm:"not null" json:"imageUrl"`
	Order    int    `gorm:"default:0" json:"order"`
	IsActive bool   `gorm:"default:true" json:"isActive"`
}

func (Backdrop) TableName() string { return "backdrops" }
