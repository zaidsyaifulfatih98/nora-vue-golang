package models

type Feature struct {
	BaseModel
	Title       string `gorm:"size:100;not null" json:"title"`
	Description string `gorm:"type:text;not null" json:"description"`
	ImageURL    string `gorm:"not null" json:"imageUrl"`
	Order       int    `gorm:"default:0" json:"order"`
	IsActive    bool   `gorm:"default:true" json:"isActive"`
}

func (Feature) TableName() string { return "features" }
