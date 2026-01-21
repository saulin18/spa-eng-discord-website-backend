package model

import "time"

type LinkReport struct {
	ID         string    `json:"id"`
	PodcastID  string    `json:"podcastId"`
	ReporterIP *string   `json:"reporterIp,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type LinkReportCount struct {
	PodcastID string `json:"podcastId"`
	Count     int    `json:"count"`
}
