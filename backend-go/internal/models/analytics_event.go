package models

// AnalyticsEvent records a single "page_view" or "cta_click" fired from the
// public site. It's intentionally schema-light (no user/session tracking) —
// just enough to answer "how many visits" and "how many CTA clicks" on the
// dashboard.
type AnalyticsEvent struct {
	BaseModel
	EventType string `gorm:"size:20;not null;index" json:"eventType"`
	Path      string `gorm:"size:255" json:"path"`
	CtaLabel  string `gorm:"size:100;index" json:"ctaLabel"`
}

func (AnalyticsEvent) TableName() string { return "analytics_events" }
