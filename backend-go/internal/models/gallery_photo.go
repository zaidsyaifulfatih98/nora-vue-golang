package models

type GalleryPhoto struct {
	BaseModel
	URL      string `gorm:"not null" json:"url"`
	Caption  string `gorm:"size:150" json:"caption"`
	Order    int    `gorm:"default:0" json:"order"`
	IsActive bool   `gorm:"default:true" json:"isActive"`
}

func (GalleryPhoto) TableName() string { return "gallery_photos" }
