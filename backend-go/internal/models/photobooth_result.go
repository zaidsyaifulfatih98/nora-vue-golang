package models

// PhotoboothResult is a guest's finished digital photobooth image, saved via
// the result screen's "Simpan" button, kept so the dashboard can browse the
// collection later. DownloadURL is derived from ImageURL (Cloudinary's
// fl_attachment flag) rather than stored separately.
type PhotoboothResult struct {
	BaseModel
	ImageURL string `gorm:"not null" json:"imageUrl"`
}

func (PhotoboothResult) TableName() string { return "photobooth_results" }
