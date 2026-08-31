package models

import (
	"database/sql/driver"
	"fmt"
)

// JSONText is a JSON-shaped value stored in a plain "text" column. Unlike
// encoding/json.RawMessage, it implements sql.Scanner/driver.Valuer, which
// pgx requires to read/write the column — RawMessage alone fails with
// "unsupported Scan ... storing driver.Value type string into *json.RawMessage".
// (gorm.io/datatypes.JSON was tried instead, but it hardcodes the Postgres
// column type to jsonb regardless of an explicit `type:text` tag, which
// broke migrating the existing text column.)
type JSONText []byte

func (j JSONText) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("[]"), nil
	}
	return j, nil
}

func (j *JSONText) UnmarshalJSON(data []byte) error {
	*j = append((*j)[:0], data...)
	return nil
}

func (j *JSONText) Scan(value any) error {
	if value == nil {
		*j = JSONText("[]")
		return nil
	}
	switch v := value.(type) {
	case string:
		*j = JSONText(v)
	case []byte:
		*j = JSONText(append([]byte(nil), v...))
	default:
		return fmt.Errorf("unsupported Scan type %T for JSONText", value)
	}
	return nil
}

func (j JSONText) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return string(j), nil
}

// PhotoboothFrame is the overlay used by the digital photobooth "try it"
// flow: a PNG with transparent cutout boxes where the guest's captured
// photos show through. This is distinct from FrameTemplate, which is a
// decorative, non-transparent design shown on the landing page.
//
// Slots holds the exact position of each transparent cutout box, as a JSON
// array of {x,y,width,height} fractions (0..1) relative to the frame image's
// own dimensions, e.g. [{"x":0.08,"y":0.04,"width":0.6,"height":0.28}, ...].
// It is authored visually in the dashboard against the uploaded frame image,
// since arbitrary frame artwork has no predictable slot layout.
type PhotoboothFrame struct {
	BaseModel
	Name     string   `gorm:"size:100;not null" json:"name"`
	ImageURL string   `gorm:"not null" json:"imageUrl"`
	Slots    JSONText `gorm:"type:text;not null;default:'[]'" json:"slots"`
	Order    int      `gorm:"default:0" json:"order"`
	IsActive bool     `gorm:"default:true" json:"isActive"`
}

func (PhotoboothFrame) TableName() string { return "photobooth_frames" }
