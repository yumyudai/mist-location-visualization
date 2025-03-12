package models

import (
	"time"
)

// Map represents a floor map in the system
type Map struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"-"`
	SiteId    string    `json:"site_id"`
	Width     int64     `json:"width"`
	Height    int64     `json:"height"`
	Ppm       float64   `json:"ppm"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Zone represents a defined area on a map
type Zone struct {
	Name      string    `json:"name"`
	Id        string    `gorm:"primaryKey;not null" json:"id"`
	MapId     string    `json:"map_id"`
	SiteId    string    `json:"site_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Entity represents a tracked device or asset in the system
type Entity struct {
	Mac         string    `gorm:"primaryKey;not null" json:"mac"`
	MapId       string    `json:"map_id"`
	Name        string    `json:"name"`
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Lastseen    float64   `json:"last_seen"`
	ZoneId      string    `json:"zone_id"`
	ZoneName    string    `json:"zone_name"`
	DisplayName string    `json:"display_name"`
	DisplayOrg  string    `json:"display_org"`
	LastRefresh time.Time `json:"refreshed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
