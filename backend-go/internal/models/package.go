package models

import (
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Package struct {
	BaseModel
	Name          string          `gorm:"size:50;not null" json:"name"`
	Price         decimal.Decimal `gorm:"type:decimal;not null" json:"price"`
	Duration      string          `gorm:"size:100;not null" json:"duration"`
	Description   string          `gorm:"type:text;not null" json:"description"`
	Features      pq.StringArray  `gorm:"type:text[]" json:"features"`
	NameEn        *string         `gorm:"size:50" json:"nameEn"`
	DurationEn    *string         `gorm:"size:100" json:"durationEn"`
	DescriptionEn *string         `gorm:"type:text" json:"descriptionEn"`
	FeaturesEn    pq.StringArray  `gorm:"type:text[]" json:"featuresEn"`
	IsPopular     bool            `gorm:"default:false" json:"isPopular"`
	IsActive      bool            `gorm:"default:true" json:"isActive"`
	Order         int             `gorm:"default:0" json:"order"`
}

func (Package) TableName() string { return "packages" }
