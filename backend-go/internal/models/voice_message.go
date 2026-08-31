package models

// VoiceMessage is a guest's recorded voice greeting, submitted from the
// digital photobooth result screen and collected for the couple/organizer
// to listen to in the dashboard. PhotoURL optionally links back to the
// digital photobooth result the guest made in the same session, for context.
type VoiceMessage struct {
	BaseModel
	GuestName string `gorm:"size:100" json:"guestName"`
	AudioURL  string `gorm:"not null" json:"audioUrl"`
	PhotoURL  string `json:"photoUrl"`
}

func (VoiceMessage) TableName() string { return "voice_messages" }
