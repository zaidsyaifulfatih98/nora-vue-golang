package models

// PhotoboothResult is a guest's finished digital photobooth image, saved via
// the result screen's "Simpan" button, kept so the dashboard can browse the
// collection later. DownloadURL is derived from ImageURL (Cloudinary's
// fl_attachment flag) rather than stored separately.
type PhotoboothResult struct {
	BaseModel
	ImageURL string `gorm:"not null" json:"imageUrl"`
	// OwnerID is nil for results saved from the main site (today's
	// behavior). When set, it scopes the result to a DIGITAL_PHOTOBOOTH
	// customer account whose sub-URL the guest used.
	OwnerID *string `gorm:"index" json:"ownerId,omitempty"`
}

func (PhotoboothResult) TableName() string { return "photobooth_results" }
