package models

type Role string

const (
	RoleSuperAdmin Role = "SUPER_ADMIN"
	RoleAdmin      Role = "ADMIN"
)

type User struct {
	BaseModel
	FirstName string `gorm:"size:25;not null" json:"firstName"`
	LastName  string `gorm:"size:25;not null" json:"lastName"`
	Email     string `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	Role      Role   `gorm:"type:text;not null;default:ADMIN" json:"role"`
}

func (User) TableName() string { return "users" }
