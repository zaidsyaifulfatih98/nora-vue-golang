package models

type FrameTemplateFit string

const (
	FrameFitCover   FrameTemplateFit = "COVER"
	FrameFitContain FrameTemplateFit = "CONTAIN"
)

type FrameTemplate struct {
	BaseModel
	Name        string           `gorm:"size:100;not null" json:"name"`
	Description string           `gorm:"type:text;not null" json:"description"`
	ImageURL    string           `gorm:"not null" json:"imageUrl"`
	Fit         FrameTemplateFit `gorm:"type:text;not null;default:COVER" json:"fit"`
	Order       int              `gorm:"default:0" json:"order"`
	IsActive    bool             `gorm:"default:true" json:"isActive"`
}

func (FrameTemplate) TableName() string { return "frame_templates" }
