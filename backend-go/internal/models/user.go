package models

type Role string

const (
	RoleSuperAdmin         Role = "SUPER_ADMIN"
	RoleAdmin              Role = "ADMIN"
	RoleDigitalPhotobooth  Role = "DIGITAL_PHOTOBOOTH"
	RoleSoftwarePhotobooth Role = "SOFTWARE_PHOTOBOOTH"
)

// IsCustomer reports whether the role is a customer account type (as
// opposed to SUPER_ADMIN/ADMIN) — DIGITAL_PHOTOBOOTH and SOFTWARE_PHOTOBOOTH
// customers each get their own slug and own owner-scoped
// PhotoboothFrame/PhotoboothResult/VoiceMessage rows.
func (r Role) IsCustomer() bool {
	return r == RoleDigitalPhotobooth || r == RoleSoftwarePhotobooth
}

// Slug is only set for customer accounts (DIGITAL_PHOTOBOOTH,
// SOFTWARE_PHOTOBOOTH) — it builds their public sub-URL
// (e.g. /digital-photobooth/{slug} or /software-photobooth/{slug}) and is
// used to scope PhotoboothFrame/PhotoboothResult/VoiceMessage rows to their
// account.
type User struct {
	BaseModel
	FirstName string  `gorm:"size:25;not null" json:"firstName"`
	LastName  string  `gorm:"size:25;not null" json:"lastName"`
	Email     string  `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password  string  `gorm:"not null" json:"-"`
	Role      Role    `gorm:"type:text;not null;default:ADMIN" json:"role"`
	Slug      *string `gorm:"uniqueIndex" json:"slug,omitempty"`
}

func (User) TableName() string { return "users" }
