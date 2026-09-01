package models

type Review struct {
	BaseModel
	Name         string  `gorm:"size:100;not null" json:"name"`
	EventLabel   string  `gorm:"size:150;not null" json:"eventLabel"`
	Quote        string  `gorm:"type:text;not null" json:"quote"`
	EventLabelEn *string `gorm:"size:150" json:"eventLabelEn"`
	QuoteEn      *string `gorm:"type:text" json:"quoteEn"`
	Rating       int     `gorm:"default:5" json:"rating"`
	IsPublished  bool    `gorm:"default:true" json:"isPublished"`
	Order        int     `gorm:"default:0" json:"order"`
}

func (Review) TableName() string { return "reviews" }
